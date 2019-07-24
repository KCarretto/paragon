package starlark

import (
	"errors"
	"fmt"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

var (
	// ErrScriptNotFound occurs when attempting to resolve a script that is not registered with the
	// flux engine.
	ErrScriptNotFound = errors.New("could not resolve script")

	// ErrMethodNotFound occurs when attempting to resolve a function that does not exist.
	ErrMethodNotFound = errors.New("could not resolve method")

	// ErrInvalidParamType occurs when a script passes an invalid parameter to a golang built-in method.
	ErrInvalidParamType = errors.New("invalid parameter type")
)

// A ScriptOption allows for flexible configuration of newly created flux scripts.
type ScriptOption func(*Script)

// Script represents starlark code that is executable by the flux engine.
type Script struct {
	name    string
	content string
	symbols starlark.StringDict
}

// NewScript returns a new flux script that is ready to be registered with the flux engine.
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

// A Processor is used by the flux engine to handle additional processing logic during script registration.
type Processor func(*Engine, *Script)

// Engine represents a starlark interpeter equipped with flux domain-specific functionality.
type Engine struct {
	globals starlark.StringDict
	scripts map[string]*Script

	processors []Processor

	// TODO: sync.pool instead of single thread
	thread *starlark.Thread
}

// RegisterScript evaluates and registers the provided script with the flux engine, returning any
// compilation errors it may have caused. If the script has already been registered, it will not be
// executed again. To update the flux engine's cache of this script, invoke UpdateScript() instead.
func (e *Engine) RegisterScript(fluxScript *Script) error {
	cached, ok := e.scripts[fluxScript.name]
	if !ok || cached == nil {
		return e.UpdateScript(fluxScript)
	}
	return nil
}

// UpdateScript evaluates and registers the provided script with the flux engine, returning any
// compilation errors it may have caused. The script will be re-evaluated even if it has been
// previously registered. In this case, flux will not replace the cached version of the script
// unless evaluation completes without errors.
func (e *Engine) UpdateScript(fluxScript *Script) error {
	// Parse script
	symbols, err := starlark.ExecFile(e.thread, fluxScript.name, fluxScript.content, e.globals)
	if err != nil {
		return err
	}

	// Update script
	fluxScript.symbols = symbols
	e.scripts[fluxScript.name] = fluxScript

	// Flux engine post-processing
	for _, processor := range e.processors {
		processor(e, fluxScript)
	}

	return nil
}

func (e *Engine) loadModule(_ *starlark.Thread, module string) (starlark.StringDict, error) {
	script, ok := e.scripts[module]
	if !ok || script == nil {
		return nil, ErrScriptNotFound
	}
	return script.symbols, nil
}

// RemoveScript removes a flux script from the engine's cache. If the script does not exist, this is
// a no op.
func (e *Engine) RemoveScript(name string) {
	delete(e.scripts, name)
}

// RemoveAllScripts removes all flux scripts from the engine's cache.
func (e *Engine) RemoveAllScripts(name string) {
	e.scripts = map[string]*Script{}
}

// Globals returns starlark values that will be used for module execution with the flux engine.
func (e *Engine) Globals() starlark.StringDict {
	return e.globals
}

// Call executes a method defined in the provided script.
func (e *Engine) Call(scriptName string, methodName string, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	script, ok := e.scripts[scriptName]
	if !ok || script == nil {
		return nil, ErrScriptNotFound
	}

	method, ok := script.symbols[methodName]
	if !ok || method == nil {
		return nil, ErrMethodNotFound
	}

	return starlark.Call(e.thread, method, args, kwargs)
}

func (e *Engine) print(_ *starlark.Thread, msg string) {
	// TODO: Configure outputter
	fmt.Println(msg)
}

func (e *Engine) newThread() *starlark.Thread {
	thread := &starlark.Thread{
		// TODO: Thread IDs?
		Name:  "flux_engine_thread",
		Print: e.print,
		Load:  e.loadModule,
	}

	return thread
}

// DefaultGlobals returns the default set of global starlark values used by the flux engine.
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

// New creates and returns a new flux engine implementation.
func New(options ...Option) *Engine {
	engine := &Engine{
		globals: DefaultGlobals(),
		scripts: map[string]*Script{},
	}
	engine.thread = engine.newThread()

	for _, opt := range options {
		opt(engine)
	}

	return engine
}
