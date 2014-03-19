// Copyright 2013 Masterminds.

package cookoo

import (
	"fmt"
)

// A Registry contains the the callback routes and the commands each
// route executes.
type Registry struct {
	routes            map[string]*routeSpec
	orderedRouteNames []string
	currentRoute      *routeSpec
}

// NewRegistry returns a new initialized registry.
func NewRegistry() *Registry {
	r := new(Registry)
	r.Init()
	return r
}

// Init initializes a registry. If a Registry is created through a means other
// than NewRegistry Init should be called on it.
func (r *Registry) Init() *Registry {
	// Why 8?
	r.routes = make(map[string]*routeSpec, 8)
	r.orderedRouteNames = make([]string, 0, 8)
	return r
}

// Route specifies a new route to add to the registry.
func (r *Registry) Route(name, description string) *Registry {

	// Create the route spec.
	route := new(routeSpec)
	route.name = name
	route.description = description
	route.commands = make([]*commandSpec, 0, 4)

	// Add the route spec.
	r.currentRoute = route
	r.routes[name] = route
	r.orderedRouteNames = append(r.orderedRouteNames, name)

	return r
}

// Does adds a command to the end of the chain of commands for the current
// (most recently specified) route.
func (r *Registry) Does(cmd Command, commandName string) *Registry {

	// Configure command spec.
	spec := new(commandSpec)
	spec.name = commandName
	spec.command = cmd

	// Add command spec.
	r.currentRoute.commands = append(r.currentRoute.commands, spec)

	return r
}

// Using specifies a paramater to use for the most recently specified command
// as set by Does.
func (r *Registry) Using(name string) *Registry {
	// Look up the last command added.
	lastCommand := r.lastCommandAdded()

	// Create a new spec.
	spec := new(paramSpec)
	spec.name = name

	// Add it to the list.
	lastCommand.parameters = append(lastCommand.parameters, spec)
	return r
}

// WithDefault specifies the default value for the most recently specified
// parameter as set by Using.
func (r *Registry) WithDefault(value interface{}) *Registry {
	param := r.lastParamAdded()
	param.defaultValue = value
	return r
}

// From sepcifies where to get the value from for the most recently specified
// paramater as set by Using.
func (r *Registry) From(fromVal string) *Registry {
	param := r.lastParamAdded()
	param.from = fromVal
	return r
}

// Done ends specifing a list of routes on a Registry. This is an optional
// method.
// TODO: Can this be removed?
func (r *Registry) Done() *Registry {
	return r
}

// Get the last parameter for the last command added.
func (r *Registry) lastParamAdded() *paramSpec {
	cspec := r.lastCommandAdded()
	last := len(cspec.parameters) - 1
	return cspec.parameters[last]
}

// Includes makes the commands from another route avaiable on this route.
func (r *Registry) Includes(route string) *Registry {

	// Not that we don't clone commands; we just add the pointer to the current
	// route.
	spec := r.routes[route]
	if spec == nil {
		panicString := fmt.Sprintf("Could not find route %s. Skipping include.", route)
		panic(panicString)
	}
	for _, cmd := range spec.commands {
		r.currentRoute.commands = append(r.currentRoute.commands, cmd)
	}
	return r
}

// RouteSpec gets a ruote cased on its name.
func (r *Registry) RouteSpec(routeName string) (spec *routeSpec, ok bool) {
	spec, ok = r.routes[routeName]
	return
}

// Routes gets an unordered map of routes names to route specs.
//
// If order is important, use RouteNames to get the names (in order).
func (r *Registry) Routes() map[string]*routeSpec {
	return r.routes
}

// RouteNames gets a slice containing the names of every registered route.
//
// The route names are returned in the order they were added to the
// registry. This is useful to some resolvers, which apply rules in order.
func (r *Registry) RouteNames() []string {
	return r.orderedRouteNames
	/*
		names := make([]string, len(r.routes))
		i := 0
		for k := range r.routes {
			names[i] = k
			i++
		}
		return names
	*/
}

// Look up the last command.
func (r *Registry) lastCommandAdded() *commandSpec {
	lastIndex := len(r.currentRoute.commands) - 1
	return r.currentRoute.commands[lastIndex]
}

type routeSpec struct {
	name, description string
	commands          []*commandSpec
}

type commandSpec struct {
	name       string
	command    Command
	parameters []*paramSpec
}

type paramSpec struct {
	name         string
	defaultValue interface{}
	from         string
}
