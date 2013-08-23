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

	names := []string{"/foo/bar/baz", "/foo/bar/blurp", "/foo/car/baz", "/foo/anything/baz", "/foo/far/baz"}
	expects := []string{"/foo/bar/baz", "/foo/bar/*", "/foo/c??/baz", "/foo/*/baz", "/foo/[cft]ar/baz"}
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
