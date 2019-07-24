package execution

import (
	"fmt"
<<<<<<< HEAD
	"testing"
=======
>>>>>>> moved files

	"go.starlark.net/starlark"
)

// TODO: Subscribe to NATS
// TODO: Dynamically load scripts (i.e. w/ events)

const loadMe string = `
<<<<<<< HEAD
def loadedFunc(loadedStr):
=======
def loadedFunc(loadedStr) -> str:
>>>>>>> moved files
	print("I am so jacked br0: "+loadedStr)
	return "5"
`

const demoScript string = `
load("loadable_script_demo", "loadedFunc")

def say_hi(name):
	print("Hello, "+name)
def count():
	nums = [str(x) for x in range(10)]
	print("\nHere's some numbers: " + ", ".join(nums))
	newVal = getSum("Here's my string with ...")
	print("ooh, shiny: " + newVal)
	x = loadedFunc("...but not really")
	print(x)
`

<<<<<<< HEAD
func mainfunc() {
=======
func main() {
>>>>>>> moved files
	loadDemo := NewScript("loadable_script_demo", loadMe)
	demo := NewScript("demo", demoScript)

	engine := New(
		WithGlobalFunction(
			"getSum",
			func(args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
				toStr, ok := starlark.AsString(args[0])
				if !ok {
					return nil, ErrInvalidParamType
				}
				fmt.Printf(
					"GOLANG: I'm running in golang, and can do anything golang can do: %s\n",
					toStr,
				)
				return starlark.String(toStr + " Extra Data!"), nil
			}))

	if err := engine.RegisterScript(loadDemo); err != nil {
		panic(err)
	}
	if err := engine.RegisterScript(demo); err != nil {
		panic(err)
	}

	if _, err := engine.Call("demo", "say_hi", starlark.Tuple{starlark.String("World")}, nil); err != nil {
		panic(err)
	}
	if _, err := engine.Call("demo", "count", nil, nil); err != nil {
		panic(err)
	}
}

<<<<<<< HEAD
func TestIt(t *testing.T) {
	mainfunc()
}

=======
>>>>>>> moved files
/*
 Each script declares events to listen for.
 Based on the events, register handlers with a map.
 Subscribe to those events via NATS.
 When an Event is received, call process(<event>)

 Example:
	def subscribe():
		return [Events.NewCreds, Events.FileModified]

	def process(event):
		print("Processing event: " + event.Topic)
*/
