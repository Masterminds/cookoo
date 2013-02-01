package cookoo

const VERSION = "0.0.1"

// Create a new Cookoo app.
func Cookoo() (reg *Registry, router *Router, cxt Context) {
	cxt = NewContext()
	reg = NewRegistry()
	router = NewRouter(reg)
	return
}

// Execute a command and return a result.
// A Cookoo app has a registry, which has zero or more routes. Each route
// executes a sequence of zero or more commands. A command is of this type.
type Command func(cxt Context, params *Params) interface{}

