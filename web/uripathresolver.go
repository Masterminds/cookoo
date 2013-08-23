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
// The most exact match "wins". E.g. for registry items `/foo/bar` and `/foo/**`, if the 
// URI path is `/foo/bar`, the `/foo/bar` entry will match first.

package web

import (
	"github.com/masterminds/cookoo"
	"path"
)

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
