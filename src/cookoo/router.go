package cookoo

const IS_ROUTER = true

type Router struct {
	registry *Registry
	resolver *RequestResolver
}

interface RequestResolver {
	func Init(registry *Registry)
	func Resolve(path string) string
}

func (r *Router) Init(registry *Registry) *Router {
	r.registry = registry
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
func (r *Router) SetRequestResolver (resolver *RequestResolver) {
}

// Resolver a given string into a route name.
func (r *Router) ResolveRequest(name string, cxt ExecutionContext) string {
}

// Do a request.
func (r *Router) HandleRequest(name string, cxt ExecutionContext, taint bool) {
}

// Check whether the given request is in the registry.
//
// This will resolve the name first.
func (r *Router) HasRequest(name string) bool {
	route := r.ResolveRequest(name)
}



