package web

import (
	"github.com/masterminds/cookoo"
	"testing"
	"fmt"
)

func TestUriPathResolver (t *testing.T) {
	reg, router, cxt := cookoo.Cookoo()
	resolver := new(URIPathResolver)
	resolver.Init(reg)

	router.SetRequestResolver(resolver)

	// ORDER IS IMPORTANT!
	reg.Route("/foo/bar/baz", "test")
	reg.Route("/foo/bar/*", "test")
	reg.Route("/foo/c??/baz", "test")
	reg.Route("/foo/[cft]ar/baz", "test")
	reg.Route("/foo/*/baz", "test")
	reg.Route("/foo/[0-9]*/baz", "test")
	reg.Route("/*/*/*", "test")

	reg.Route("GET /foo/bar/baz", "Test with verb")
	reg.Route("POST /foo/bar/baz", "Test with verb")
	reg.Route("DELETE /foo/bar/baz", "Test with verb")
	reg.Route("* /foo/bar/baz", "Test with verb")
	reg.Route("* /foo/last", "Test with verb")

	names := []string{"/foo/bar/baz", "/foo/bar/blurp", "/foo/car/baz", "/foo/anything/baz", "/foo/far/baz", "POST /foo/bar/baz", "GET /foo/last"}
	expects := []string{"/foo/bar/baz", "/foo/bar/*", "/foo/c??/baz", "/foo/*/baz", "/foo/[cft]ar/baz", "POST /foo/bar/baz", "* /foo/last"}
	for i, name := range names {
		resolved, err := router.ResolveRequest(name, cxt)
		if err != nil {
			t.Error("Unexpected resolver error", err)
		}
		if resolved != expects[i] {
			t.Error(fmt.Sprintf("! Expected to find %s at %d; found %s", expects[i], i, resolved))
		}
	}
}
