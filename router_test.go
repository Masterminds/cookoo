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

func MockCommand(cxt Context, params *Params) (interface{}, Interrupt) {
	//println("Mock command")
	return true, nil
}

func RerouteCommand(cxt Context, params *Params) (interface{}, Interrupt) {
	route := params.Get("route", "default").(string)
	return nil, &Reroute{route}
}

func FetchParams(cxt Context, params *Params) (interface{}, Interrupt) {
	return params, nil;
}

func RecoverableErrorCommand(cxt Context, params *Params) (interface{}, Interrupt) {
	return nil, &RecoverableError{"Blarg"}
}

func FatalErrorCommand(cxt Context, params *Params) (interface{}, Interrupt) {
	return nil, &FatalError{"Blarg"}
}

type MockDatasource struct {
	RetVal string;
}

func (ds *MockDatasource) Value(key string) interface{} {
	return ds.RetVal;
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

func TestFromValues(t *testing.T) {
	reg, router, cxt := Cookoo()

	cxt.Add("test1", 1234)
	cxt.AddDatasource("test2", "foo")

	ds := new(MockDatasource);
	ds.RetVal = "1234"
	cxt.AddDatasource("foo", ds);

	reg.
		Route("mock", "Test from.").
			Does(FetchParams, "first").
				Using("test1").From("cxt:test1").
				Using("test2").From("datasource:test2").
				Using("test3").From("foo:test3").
				Using("test4").WithDefault("test4").From("NONE:none").
				Using("test5").WithDefault("Z").From("NONE:none foo:test3 cxt:test1").
				Using("test6").From("None:none")

		e := router.HandleRequest("mock", cxt, true);
		if e != nil {
			t.Error("Unexpected: ", e.Error());
		}

		params, ok := cxt.Get("first").(*Params);
		if !ok {
			t.Error("! Expected a Params object.")
		}

		test1, ok := params.Has("test1");
		if !ok {
			t.Error("! Expected a value in cxt:test1");
		}
		if test1.(int) != 1234 {
			t.Error("! Expected test1 to return 1234. Got ", test1);
		}


		test2, ok := params.Has("test2");
		if !ok {
			t.Error("! Expected a value in cxt:test1");
		}
		if test2.(string) != "foo" {
			t.Error("! Expected test2 to return 'foo'. Got ", test2);
		}

		test3, ok := params.Has("test3");
		if !ok {
			t.Error("! Expected default value");
		}
		if test3.(string) != "1234" {
			t.Error("! Expected test4 to return '1234'. Got ", test3);
		}

		test4, ok := params.Has("test4");
		if !ok {
			t.Error("! Expected default value");
		}
		if test4.(string) != "test4" {
			t.Error("! Expected test4 to return 'test4'. Got ", test4);
		}

		// We expect that in this case the first match in the From clause
		// will be returned, which is the value of foo:test3.
		test5, ok := params.Has("test3");
		if !ok {
			t.Error("! Expected default value");
		}
		if test5.(string) != "1234" {
			t.Error("! Expected test5 to return '1234'. Got ", test5);
		}

		param, ok := params.Has("test6");
		if !ok {
			t.Error("! Expected a *Param with a nil value");
		}
		if param != nil {
			t.Error("! Expected nil value");
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

func TestReroute(t *testing.T) {
	reg, router, context := Cookoo()
	reg.
	  Route("TEST", "A test route").Does(RerouteCommand, "fake").
	  Using("route").WithDefault("TEST2").
	  Route("TEST2", "Tainted route").Does(FetchParams, "fake2").Using("foo").WithDefault("bar")
	e := router.HandleRequest("TEST", context, false)
	if e != nil {
		t.Error("! Unexpected error executing TEST")
	}

	p := context.Get("fake2")
	if p == nil {
		t.Error("! Expected data in fake2.")
	}
}

func TestRecoverableError(t *testing.T) {
	reg, router, context := Cookoo()
	reg.
	  Route("TEST", "A test route").
	  Does(RecoverableErrorCommand, "fake").
	  Does(FetchParams, "fake2").Using("foo").WithDefault("bar")
	
	e := router.HandleRequest("TEST", context, false)
	if e != nil {
		t.Error("! Unexpected error executing TEST")
	}

	p := context.Get("fake2")
	if p == nil {
		t.Error("! Expected data in fake2.")
	}
}

func TestFatalError(t *testing.T) {
	reg, router, context := Cookoo()
	reg.
	  Route("TEST", "A test route").
	  Does(FatalErrorCommand, "fake").
	  Does(FetchParams, "fake2").Using("foo").WithDefault("bar")
	
	e := router.HandleRequest("TEST", context, false)
	if e == nil {
		t.Error("! Expected error executing TEST")
	}

	p := context.Get("fake2")
	if p != nil {
		t.Error("! Expected fake2 to not get executed.")
	}
}
