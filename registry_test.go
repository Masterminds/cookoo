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

	if rspec.RouteName != "foo" {
		t.Error("! Expected route to be named 'foo'")
	}
	if rspec.Help != "A test route" {
		t.Error("! Expected description to be 'A test route'")
	}

	if len(rspec.Does) != 1 {
		t.Error("! Expected exactly one command. Found ", len(rspec.Does))
	}

	cmd := rspec.Does[0]
	if "fakeCommand" != cmd.Name {
		t.Error("! Expected to find fakeCommand command.")
	}

	if len(cmd.Parameters) != 1 {
		t.Error("! Expected exactly one paramter. Found ", len(cmd.Parameters))
	}

	pspec := cmd.Parameters[0]
	if pspec.Name != "param" {
		t.Error("! Expected the first param to be 'param'")
	}

	if pspec.DefaultValue != "value" {
		t.Error("! Expected the value to be 'value'")
	}
	fakeCxt := new(ExecutionContext)
	fakeParams := NewParamsWithValues(map[string]interface{}{"foo": "bar", "baz": 2})
	rr, err := cmd.Command(fakeCxt, fakeParams)

	if err != nil {
		t.Error("! Expected no errors.")
	}

	cRet := rr.(*FooType)

	if cRet.test != 5 {
		t.Error("! Expected 'test' to be 5")
	}

}

func TestRouteIncludes(t *testing.T) {
	reg := new(Registry)
	reg.Init()

	reg.Route("foo", "A test route").
		Does(AnotherCommand, "fakeCommand").
		Using("param").WithDefault("foo").
		Route("bar", "Another test route").
		Does(AnotherCommand, "fakeCommand2").
		Using("param").WithDefault("bar").
		Includes("foo").
		Does(AnotherCommand, "fakeCommand3").
		Using("param").WithDefault("baz")

	expecting := []string{"fakeCommand2", "fakeCommand", "fakeCommand3"}
	spec, ok := reg.RouteSpec("bar")
	if !ok {
		t.Error("! Expected to find a route named 'bar'")
	}
	for i, k := range expecting {
		if k != spec.Does[i].Name {
			t.Error(fmt.Sprintf("Expecting %s at position %d; got %s", k, i, spec.Does[i].Name))
		}
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("! Failed to panic when including commands for a route that does not exist.")
		}
	}()
	reg2 := new(Registry)
	reg2.Init()

	reg2.Route("foo", "A test route").
		Includes("bar")
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

	if spec.RouteName != "foo" {
		t.Error("! Expected a spec named 'foo'")
	}

	param := spec.Does[0].Parameters[1]
	if v, ok := param.DefaultValue.(Context); !ok {
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

func TestNewStyleRoutes(t *testing.T) {
	reg := NewRegistry()
	route := &RouteSpec{
		RouteName: "test",
		Help:      "This is a test",
		Does: []*CommandSpec{
			&CommandSpec{
				Name:    "foo",
				Command: AnotherCommand,
				Parameters: []*ParamSpec{
					&ParamSpec{
						Name:         "one",
						DefaultValue: "two",
						From:         "cxt:three",
					},
				},
			},
		},
	}

	reg.AddRoute(route)
	spec, ok := reg.RouteSpec("test")
	if !ok || spec == nil {
		t.Errorf("Expected to find the test route.")
	}

	if len(route.Does) != 1 {
		t.Errorf("Expected one command, got %d", len(route.Does))
	}
	if route.Does[0].Name != "foo" {
		t.Errorf("Expected command named 'foo', got '%s'", route.Does[0].Name)
	}

}
