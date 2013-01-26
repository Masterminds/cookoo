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

func TestParseFromStatement(t *testing.T) {
	str := "foo:bar foo:baz blarg:urg"
	res := parseFromStatement(str)
	if len(res) != 3 {
		t.Error("! Expected length 3, got ", len(res))
	}
	exp := res[0]
	if exp.source != "foo" {
		t.Error("! Expected foo, got ", exp.source)
	}
	if exp.key != "bar" {
		t.Error("! Expected bar, got ", exp.source)
	}

	exp = res[1]
	if exp.source != "foo" {
		t.Error("! Expected foo, got ", exp.source)
	}
	if exp.key != "baz" {
		t.Error("! Expected baz, got ", exp.source)
	}

	exp = res[2]
	if exp.source != "blarg" {
		t.Error("! Expected blarg, got ", exp.source)
	}
	if exp.key != "urg" {
		t.Error("! Expected urg, got ", exp.source)
	}
}

func TestParseFromVal(t *testing.T) {
	fr := "test:foo"

	r := parseFromVal(fr)
	name := r.source
	val := r.key
	if name != "test" {
		t.Error("Expected 'test', got ", name)
	}
	if val != "foo" {
		t.Error("Expected 'foo', got ", val)
	}

	fr = "test"
	r = parseFromVal(fr)
	name = r.source
	val = r.key
	if name != "test" {
		t.Error("Expected 'test', got ", name)
	}
	if val != "" {
		t.Error("Expected an empty string, got ", val)
	}

	fr = "test:"
	r = parseFromVal(fr)
	name = r.source
	val = r.key
	if name != "test" {
		t.Error("Expected 'test', got ", name)
	}
	if val != "" {
		t.Error("Expected an empty string, got ", val)
	}

	fr = "test:foo:bar:baz"
	r = parseFromVal(fr)
	name = r.source
	val = r.key
	if name != "test" {
		t.Error("Expected 'test', got ", name)
	}
	if val != "foo:bar:baz" {
		t.Error("Expected 'foo:bar:baz' string, got ", val)
	}

	fr = ""
	r = parseFromVal(fr)
	name = r.source
	val = r.key
	if name != "" {
		t.Error("Expected empty string, got ", name)
	}
	if val != "" {
		t.Error("Expected an empty string string, got ", val)
	}
}

func TestHandleRequest(t *testing.T) {
	reg, router, context := Cookoo()
	reg.
	  Route("TEST", "A test route").Does(MockCommand, "fake").
	  Route("@tainted", "Tainted route").Does(MockCommand, "fake2").
		Route("Several", "Test multiple.").
			Does(MockCommand, "first").
			Does(MockCommand, "second").
			Does(MockCommand, "third")

	e := router.HandleRequest("TEST", context, true)
	if e != nil {
		t.Error("Unexpected: ", e.Error());
	}

	e = router.HandleRequest("@tainted", context, true)
	if e == nil {
		t.Error("Expected tainted route to not run protected name.");
	}

	e = router.HandleRequest("@tainted", context, false)
	if e != nil {
		t.Error("Unexpected: ", e.Error());
	}

	router.HandleRequest("NO Such Route", context, false)

	context = NewContext()
	router.HandleRequest("Several", context, false)
	if context.Len() != 3 {
		t.Error("! Expected three items in the context, got ", context.Len())
	}
}
