package cookoo

import (
	"fmt"
	"testing"
)

func ExampleCookoo() {
	reg, router, cxt := Cookoo()
	reg.AddRoute(Route{
		Name: "hello",
		Help: "Logs 'hello world' to standard output",
		Does: Tasks{
			Cmd{
				Name: "log",
				// Usually we define functions elsewhere so we can re-use them.
				Fn: func(c Context, p *Params) (interface{}, Interrupt) {
					fmt.Println("Hello World")
					return nil, nil
				},
				Using: []Param{
					Param{
						Name:         "msg",
						DefaultValue: "Hello World",
					},
				},
			},
		},
	})

	router.HandleRequest("hello", cxt, false)
	// Output:
	// Hello World
}

func TestCookooForCoCo(t *testing.T) {
	registry, router, cxt := Cookoo()

	cxt.Put("Answer", 42)

	lifeUniverseEverything := cxt.Get("Answer", nil)

	if lifeUniverseEverything != 42 {
		t.Error("! Context is not working.")
	}

	registry.Route("foo", "test")

	ok := router.HasRoute("foo")

	if !ok {
		t.Error("! Router does not have 'foo' route.")
	}
}
