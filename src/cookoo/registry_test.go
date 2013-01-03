package cookoo

import (
	"testing"
//	"registry"
	"fmt"
)

/*
type FakeCommand struct {
	Command
}
*/

func FakeCommand(cxt *ExecutionContext, params map[string]*interface{}) bool {
	fmt.Println("Got here")

	return true
}

func TestBasicRoute (t *testing.T) {
	reg := new(Registry)
	reg.Init()

	reg.Route("foo", "A test route")
	reg.Does(FakeCommand, "fakeCommand").Using("param").WithDefaultValue("value")

	// Now do something to test.
	routes := reg.Routes()

	if len(routes) != 1 {
		t.Error("! Expected one route.")
	}

	rspec := routes["foo"]

	if rspec.name != "foo" {
		t.Error("! Expected route to be named 'foo'")
	}
	if rspec.description != "A test route" {
		t.Error("! Expected description to be 'A test route'")
	}

	if len(rspec.commands) != 1 {
		t.Error("! Expected exactly one command.")
	}


	cmd := rspec.commands[0]
	fmt.Println(cmd)
	if "fakeCommand" != cmd.name {
		t.Error("! Expected to find fakeCommand command.")
	}

	if len(cmd.parameters) != 1 {
		t.Error("! Expected exactly one paramter.")
	}

	pspec := cmd.parameters[0]
	if pspec.name != "param" {
		t.Error("! Expected the first param to be 'param'")
	}
	/*
	if pspec.defaultValue.(string) != "value" {
		t.Error("! Expected the value to be 'value'")
	}
	*/
	fakeCommand := new(ExecutionContext)
	fakeParams := make(map[string]*interface{}, 2)
	//fakeParams["foo"] = "bar"
	//fakeParams["baz"] = 2
	cmd.command(fakeCommand, fakeParams)

}
