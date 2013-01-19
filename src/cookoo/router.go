package cookoo

import (
	"fmt"
)

// The request resolver.
// A request resolver is responsible for transforming a request name to
// a route name. For example, a web-specific resolver may take a URI
// and return a route name. Or it make take an HTTP verb and return a
// route name.
type RequestResolver interface {
	Init(registry *Registry)
	Resolve(path string, cxt Context) string
}

// The Cookoo router.
// A Cookoo app works by passing a request into a router, and
// relying on the router to execute the appropriate chain of
// commands.
type Router struct {
	registry *Registry
	resolver RequestResolver
}

// A basic resolver that assumes that the given request name
// *is* the route name.
type BasicRequestResolver struct {
	registry *Registry
	resolver RequestResolver
}

func NewRouter(reg *Registry) *Router {
	router := new(Router)
	router.Init(reg)
	return router
}

func (r *BasicRequestResolver) Init(registry *Registry) {
	r.registry = registry
}
func (r *BasicRequestResolver) Resolve(path string, cxt Context) string {
	return path
}

func (r *Router) Init(registry *Registry) *Router {
	r.registry = registry
	r.resolver = new(BasicRequestResolver)
	r.resolver.Init(registry)
	return r
}

// Set the registry.
func (r *Router) SetRegistry(reg *Registry) {
	r.registry = reg
}

// Set the request resolver.
// The resolver is responsible for taking an arbitrary string and
// resolving it to a registry route.
//
// Example: Take a URI and translate it to a route.
func (r *Router) SetRequestResolver (resolver RequestResolver) {
	r.resolver = resolver
}

// Get the request resolver.
func (r *Router) RequestResolver() RequestResolver {
	return r.resolver
}

// Resolve a given string into a route name.
func (r *Router) ResolveRequest(name string, cxt Context) string {
	routeName := r.resolver.Resolve(name, cxt)

	return routeName
}

// Do a request.
// This executes a request "named" name (this string is passed through the
// request resolver.) The context is cloned (shallow copy) and passed in as the
// base context.
//
// If taint is `true`, then no routes that begin with `@` can be executed. Taint
// should be set to true on anything that relies on a name supplied by an
// external client.
//
// This will do the following:
// - resolve the request name into a route name (using a RequestResolver)
// - look up the route
// - execute each command on the route in order
//
// No data is returned from a route.
func (r *Router) HandleRequest(name string, cxt Context, taint bool) {
	baseCxt := cxt.Copy()
	routeName := r.ResolveRequest(name, baseCxt)

	//go r.runRoute(routeName, cxt, taint)
	r.runRoute(routeName, cxt, taint)

}

// This checks whether or not the route exists.
// Note that this does NOT resolve a request name into a route name. This
// expects a route name.
func (r *Router) HasRoute(name string) bool {
	_, ok := r.registry.RouteSpec(name)
	return ok
}

// PRIVATE ==========================================================

func (r *Router) runRoute(route string, cxt Context, taint bool) (ok bool, err error ) {
	if len(route) == 0 {
		return true, &RouteError{"Empty route name."}
	}
	if taint && route[0] == '@' {
		return true, &RouteError{"Route is tainted. Refusing to run."}
	}
	spec, ok := r.registry.RouteSpec(route)
	if (!ok) {
		return true, &RouteError{"Route does not exist."}
	}
	fmt.Printf("Running route %s: %s\n", spec.name, spec.description)
	for i, cmd := range spec.commands {
		fmt.Printf("Command %d is %s (%T)\n", i, cmd.name, cmd.command)
		r.doCommand(cmd, cxt)
	}
	return false, nil
}

func (r *Router) doCommand(cmd *commandSpec, cxt Context) bool {
	return false
}

// Indicates that a route cannot be executed successfully.
type RouteError struct {
	Message string
}
func (e *RouteError) Error () string {
	return e.Message
}
