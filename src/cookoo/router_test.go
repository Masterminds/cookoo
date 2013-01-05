package cookoo

import (
	"testing"
)

type FakeRequestResolver struct {
	BasicRequestResolver
}

func (self *FakeRequestResolver) Resolve(name string, cxt *ExecutionContext) string {
	return "FOO"
}

func TestSetResolver (t *testing.T) {
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
