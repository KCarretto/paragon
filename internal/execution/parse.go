package execution

import (
	"errors"
	"fmt"
	"io"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

var (
	// ErrScriptNotFound occurs when attempting to resolve a script that is not registered with the
	// flux interpreter.
	ErrScriptNotFound = errors.New("could not resolve script")

	// ErrMethodNotFound occurs when attempting to resolve a function that does not exist.
	ErrMethodNotFound = errors.New("could not resolve method")

	// ErrInvalidParamType occurs when a script passes an invalid parameter to a golang built-in method.
	ErrInvalidParamType = errors.New("invalid parameter type")
)

// A ScriptOption allows for flexible configuration of newly created flux scripts.
type ScriptOption func(*Script)

// Script represents starlark code that is executable by the flux interpreter.
type Script struct {
	name    string
	content string
	symbols starlark.StringDict
}

// NewScript returns a new flux script that is ready to be registered with the flux interpreter.
func NewScript(name string, src string, options ...ScriptOption) *Script {
	script := &Script{
		name:    name,
		content: src,
	}

	for _, opt := range options {
		opt(script)
	}

	return script
}

// A Processor is used by the flux interpreter to handle additional processing logic during script registration.
type Processor func(*Interpreter, *Script)

// Interpreter represents a starlark interpeter equipped with flux domain-specific functionality.
type Interpreter struct {
	globals starlark.StringDict
	scripts map[string]*Script

	processors []Processor

	// TODO: sync.pool instead of single thread
	thread *starlark.Thread
}

// RegisterScript evaluates and registers the provided script with the flux interpreter, returning any
// compilation errors it may have caused. If the script has already been registered, it will not be
// executed again. To update the flux interpreter's cache of this script, invoke UpdateScript() instead.
func (i *Interpreter) RegisterScript(fluxScript *Script) error {
	cached, ok := i.scripts[fluxScript.name]
	if !ok || cached == nil {
		return i.UpdateScript(fluxScript)
	}
	return nil
}

// UpdateScript evaluates and registers the provided script with the flux interpreter, returning any
// compilation errors it may have caused. The script will be re-evaluated even if it has been
// previously registered. In this case, flux will not replace the cached version of the script
// unless evaluation completes without errors.
func (i *Interpreter) UpdateScript(fluxScript *Script) error {
	// Parse script
	symbols, err := starlark.ExecFile(i.thread, fluxScript.name, fluxScript.content, i.globals)
	if err != nil {
		return err
	}

	// Update script
	fluxScript.symbols = symbols
	i.scripts[fluxScript.name] = fluxScript

	// Flux interpreter post-processing
	for _, processor := range i.processors {
		processor(i, fluxScript)
	}

	return nil
}

// ExecScript executes an entire script already saved in the interpreter
func (i *Interpreter) ExecScript(name string) error {
	script, ok := i.scripts[name]
	if !ok || script != nil {
		return ErrScriptNotFound
	}
	stdout := &StdOutCapture{}
	thread := i.newThread(stdout)
	_, err := starlark.ExecFile(thread, script.name, script.content, i.globals)
	if err != nil {
		return err
	}
	stdout.Print()
	return nil
}

// loadModule is used as load within starlark
func (i *Interpreter) loadModule(_ *starlark.Thread, module string) (starlark.StringDict, error) {
	script, ok := i.scripts[module]
	if !ok || script == nil {
		return nil, ErrScriptNotFound
	}
	return script.symbols, nil
}

// RemoveScript removes a flux script from the interpreter's cache. If the script does not exist, this is
// a no op.
func (i *Interpreter) RemoveScript(name string) {
	delete(i.scripts, name)
}

// RemoveAllScripts removes all flux scripts from the interpreter's cache.
func (i *Interpreter) RemoveAllScripts(name string) {
	i.scripts = map[string]*Script{}
}

// Globals returns starlark values that will be used for module execution with the flux interpreter.
func (i *Interpreter) Globals() starlark.StringDict {
	return i.globals
}

// Call executes a method defined in the provided script.
func (i *Interpreter) Call(scriptName string, methodName string, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	script, ok := i.scripts[scriptName]
	if !ok || script == nil {
		return nil, ErrScriptNotFound
	}

	method, ok := script.symbols[methodName]
	if !ok || method == nil {
		return nil, ErrMethodNotFound
	}

	stdout := &StdOutCapture{}
	thread := i.newThread(stdout)
	val, err := starlark.Call(thread, method, args, kwargs)
	stdout.Print()
	return val, err
}

type StdOutCapture struct {
	stdout []byte
}

func (s *StdOutCapture) Write(p []byte) (int, error) {
	if s.stdout == nil {
		s.stdout = []byte{}
	}
	s.stdout = append(s.stdout, p...)
	return 0, nil // todo let's make this int real

}

func (s *StdOutCapture) Print() {
	fmt.Printf("%s", s.stdout)
}

func (i *Interpreter) print(writer io.Writer) printFunc {
	// TODO: Configure outputter
	return func(_ *starlark.Thread, msg string) {
		bytesMsg := []byte(msg)
		bytesMsg = append(bytesMsg, []byte("\n")...)
		writer.Write(bytesMsg)
	}
}

type printFunc func(*starlark.Thread, string)

func (i *Interpreter) newThread(printWriter io.Writer) *starlark.Thread {
	thread := &starlark.Thread{
		// TODO: Thread IDs?
		Name:  "flux_interpreter_thread",
		Print: i.print(printWriter),
		Load:  i.loadModule,
	}

	return thread
}

// DefaultGlobals returns the default set of global starlark values used by the flux interpreter.
func DefaultGlobals() starlark.StringDict {
	event := starlarkstruct.FromStringDict(
		starlarkstruct.Default,
		starlark.StringDict{
			"Type":  starlark.String("type"),
			"Topic": starlark.String("topic"),
			"Data":  &starlark.Dict{},
		})

	return starlark.StringDict{
		"event": event,
	}
}

// New creates and returns a new flux interpreter implementation.
func New(options ...Option) *Interpreter {
	interpreter := &Interpreter{
		globals: DefaultGlobals(),
		scripts: map[string]*Script{},
	}
	stdout := &StdOutCapture{}
	interpreter.thread = interpreter.newThread(stdout)

	for _, opt := range options {
		opt(interpreter)
	}

	return interpreter
}

// An Option for configuring a new flux interpreter.
type Option func(*Interpreter)

type globalFn = func(args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)

// WithGlobalFunction adds the provided handler to the global scope for each flux script.
func WithGlobalFunction(name string, handler globalFn) Option {
	builtIn := starlark.NewBuiltin(
		name,
		func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			return handler(args, kwargs)
		},
	)
	return func(i *Interpreter) {
		i.globals[name] = builtIn
	}
}
