package cookoo

import (
	"testing"
	//	"registry"
	"fmt"
)

type FooType struct {
	test int
}

func FakeCommand(cxt Context, params Params) (interface{}, Interrupt) {
	fmt.Println("Got here")

	var ret bool = true

	p := &ret

	return p, nil
}

func AnotherCommand(cxt Context, params *Params) (interface{}, Interrupt) {
	//ret := func() bool {return true;}
	ret := new(FooType)
	ret.test = 5

	return ret, nil
}

func TestBasicRoute(t *testing.T) {
	reg := new(Registry)
	reg.Init()

	reg.Route("foo", "A test route")
	reg.Does(AnotherCommand, "fakeCommand").Using("param").WithDefault("value")

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

	if pspec.defaultValue != "value" {
		t.Error("! Expected the value to be 'value'")
	}
	fakeCxt := new(ExecutionContext)
	fakeParams := NewParamsWithValues(map[string]interface{}{"foo": "bar", "baz": 2})
	rr, err := cmd.command(fakeCxt, fakeParams)

	if err != nil {
		t.Error("! Expected no errors.")
	}

	cRet := rr.(*FooType)

	if cRet.test != 5 {
		t.Error("! Expected 'test' to be 5")
	}

}

func TestRouteSpec(t *testing.T) {
	reg := new(Registry)
	reg.Init()

	reg.Route("foo", "A test route").
		Does(AnotherCommand, "fakeCommand").
		Using("param").WithDefault("value").
		Using("something").WithDefault(NewContext())

	spec, ok := reg.RouteSpec("foo")

	if !ok {
		t.Error("! Expected to find a route named 'foo'")
	}

	if spec.name != "foo" {
		t.Error("! Expected a spec named 'foo'")
	}

	param := spec.commands[0].parameters[1]
	if v, ok := param.defaultValue.(Context); !ok {
		t.Error("! Expected an execution context.")
	} else {
		// Canary
		v.Put("test", "test")
	}
}

func TestRouteNames(t *testing.T) {
	reg := new(Registry)
	reg.Init()
	reg.Route("one", "A route").Does(AnotherCommand, "fake")
	reg.Route("two", "A route").Does(AnotherCommand, "fake")
	reg.Route("three", "A route").Does(AnotherCommand, "fake")
	reg.Route("four", "A route").Does(AnotherCommand, "fake")
	reg.Route("five", "A route").Does(AnotherCommand, "fake")

	names := reg.RouteNames()

	if len(names) != 5 {
		t.Error("! Expected five routes, found ", len(names))
	}

	expecting := []string{"one", "two", "three", "four", "five"}
	for i, k := range expecting {
		if k != names[i] {
			t.Error(fmt.Sprintf("Expecting %s at position %d; got %s", k, i, names[i]))
		}
	}

}
