package cli

import(
	"testing"
	"flag"
	"github.com/masterminds/cookoo"
)

func Nothing(cxt cookoo.Context, params *cookoo.Params) (res interface{}, i cookoo.Interrupt) {
	return true, nil
}

func TestResolvingSimpleRoute(t *testing.T) {
	registry, router, context := cookoo.Cookoo()

	resolver := new(RequestResolver)
	resolver.Init(registry)

	router.SetRequestResolver(resolver)

	registry.Route("test", "A simple test").Does(Nothing, "nada")

	e := router.HandleRequest("test", context, false)

	if e != nil {
		t.Error("! Failed 'test' route.")
	}
}

func TestResolvingWithFlags(t *testing.T) {
	registry, router, context := cookoo.Cookoo()

	resolver := new(RequestResolver)
	resolver.Init(registry)

	flagset := flag.NewFlagSet("test", flag.ExitOnError)
	context.Add("globalFlags", flagset)

	router.SetRequestResolver(resolver)

	registry.Route("test", "Test flag parsing.").Does(Nothing, "nada")

	e := router.HandleRequest("test", context, false)
	if e != nil {
		t.Error("! Failed 'test' route.")
	}
}
