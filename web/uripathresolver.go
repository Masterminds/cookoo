package web

import (
	"github.com/masterminds/cookoo"
	"path"
)

// Resolver for transforming a URI path into a route.
//
// This is a more sophisticated path resolver, aware of
// heirarchyand wildcards.
//
// Examples:
// - URI path `/foo` matches the entry `/foo`
// - URI path `/foo/bar` could match entries like `/foo/*`, `/foo/**`, and `/foo/bar`
// - URI path `/foo/bar/baz` could match `/foo/*/baz` and `/foo/**`
//
// HTTP Verbs:
// This resolver also allows you to specify verbs at the beginning of a path:
// - "GET /foo" and "POST /foo" are separate (but legal) paths. "* /foo" will allow any verb.
// - There are no constrainst on verb name. Thus, verbs like WebDAV's PROPSET are fine, too. Or you can
//   make up your own.
//
// For verbs to work with the router, you need to configure your router to support prepending the
// verb to the route name.
//
// The most exact match "wins". E.g. for registry items `/foo/bar` and `/foo/**`, if the 
type URIPathResolver struct {
	registry *cookoo.Registry
}

func (r *URIPathResolver) Init(registry *cookoo.Registry) {
	r.registry = registry
}

// Resolve a path name based using path patterns.
//
// This resolver is designed to match path-like strings to path patterns. For example,
// the path `/foo/bar/baz` may match routes like `/foo/*/baz` or `/foo/bar/*`
func (r *URIPathResolver) Resolve(pathName string, cxt cookoo.Context) (string, error) {
	// HTTP verb support naturally falls out of the fact that spaces in paths are legal in UNIXy systems, while
	// illegal in URI paths. So presently we do no special handling for verbs. Yay for simplicity.
	for _, pattern := range r.registry.RouteNames() {
		if ok, err := path.Match(pattern, pathName); ok && err == nil {
			return pattern, nil
		} else if err != nil {
			// Bad pattern
			return pathName, err
		}
	}
	return pathName, &cookoo.RouteError{"Could not resolve route " + pathName}
}
