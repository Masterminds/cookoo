package cookoo

import (
	"testing"
//	"registry"
	"fmt"
)

type FooType struct {
	test int
}

func FakeCommand(cxt *ExecutionContext, params map[string]*interface{}) interface{} {
	fmt.Println("Got here")

	var ret bool = true

	p := &ret

	return p
}

func AnotherCommand(cxt *ExecutionContext, params map[string]*interface{}) interface{} {
	//ret := func() bool {return true;}
	ret := new(FooType)
	ret.test = 5

	return ret
}

func TestBasicRoute (t *testing.T) {
	reg := new(Registry)
	reg.Init()

	reg.Route("foo", "A test route")
	reg.Does(AnotherCommand, "fakeCommand").Using("param").WithDefaultValue("value")
	//reg.Does(FakeCommand, "fakeCommand").Using("param").WithDefaultValue("value")

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
		t.Error("! Expected exactly one command. Found ", len(rspec.commands))
	}


	cmd := rspec.commands[0]
	if "fakeCommand" != cmd.name {
		t.Error("! Expected to find fakeCommand command.")
	}

	if len(cmd.parameters) != 1 {
		t.Error("! Expected exactly one paramter. Found ", len(cmd.parameters))
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
	fakeCxt:= new(ExecutionContext)
	fakeParams := make(map[string]*interface{}, 2)
	//fakeParams["foo"] = "bar"
	//fakeParams["baz"] = 2
	var cRet *FooType = cmd.command(fakeCxt, fakeParams).(*FooType)

	if cRet.test != 5 {
		t.Error("! Expected 'test' to be 5")
	}

}

func TestRouteSpec(t *testing.T) {
	reg := new(Registry)
	reg.Init()

	reg.Route("foo", "A test route")
	reg.Does(AnotherCommand, "fakeCommand").Using("param").WithDefaultValue("value")

	spec, ok := reg.RouteSpec("foo")

	if !ok {
		t.Error("! Expected to find a route named 'foo'")
	}

	if spec.name != "foo" {
		t.Error("! Expected a spec named 'foo'")
	}
}
