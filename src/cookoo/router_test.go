package cookoo

import (
	"testing"
)

// Mock resolver
type FakeRequestResolver struct {
	BasicRequestResolver
}
// Always returns FOO.
func (self *FakeRequestResolver) Resolve(name string, cxt Context) string {
	return "FOO"
}

// Test the resolver.
func TestResolver (t *testing.T) {
	fakeCxt := new(ExecutionContext)
	registry := new(Registry)
	r := new(Router)
	r.Init(registry)

	// Canary: Check that resolver is working.
	if a := r.ResolveRequest("test", fakeCxt); a != "test" {
		t.Error("Expected path to be 'test'")
	}

	// Set and get a resolver.
	resolver := new(FakeRequestResolver)
	r.SetRequestResolver(resolver)
	resolver, ok := r.RequestResolver().(*FakeRequestResolver)

	if !ok {
		t.Error("! Resolver is not a FakeRequestResolver.")
	}

	// Make sure the new resolver works.
	path := r.ResolveRequest("test", fakeCxt)

	if path != "FOO" {
		t.Error("Expected path to be 'test'")
	}
}

func MockCommand(cxt Context, params Params) interface{} {
	println("Mock command")
	return true
}

func TestHandleRequest(t *testing.T) {
	reg, router, context := Cookoo()
	reg.
	  Route("TEST", "A test route").Does(MockCommand, "fake").
	  Route("@tainted", "Tainted route").Does(MockCommand, "fake2")

	router.HandleRequest("TEST", context, true)
	router.HandleRequest("@tainted", context, true)
	router.HandleRequest("@tainted", context, false)
	router.HandleRequest("NO Such Route", context, false)
}
