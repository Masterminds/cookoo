package cookoo

// The request resolver.
// A request resolver is responsible for transforming a request name to
// a route name. For example, a web-specific resolver may take a URI
// and return a route name. Or it make take an HTTP verb and return a
// route name.
type RequestResolver interface {
	Init(registry *Registry)
	Resolve(path string, cxt *ExecutionContext) string
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

func (r *BasicRequestResolver) Init(registry *Registry) {
	r.registry = registry
}
func (r *BasicRequestResolver) Resolve(path string, cxt *ExecutionContext) string {
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
func (r *Router) ResolveRequest(name string, cxt *ExecutionContext) string {
	routeName := r.resolver.Resolve(name, cxt)

	return routeName
}

// Do a request.
func (r *Router) HandleRequest(name string, cxt *ExecutionContext, taint bool) {
}

// Check whether the given request is in the registry.
//
// This will resolve the name first.
func (r *Router) HasRoute(name string) bool {
	return false
}



