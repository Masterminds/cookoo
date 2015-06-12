package cookoo

// Copyright 2013 Masterminds.

import (
	"fmt"
	"strings"
)

// A Registry contains the the callback routes and the commands each
// route executes.
type Registry struct {
	routes            map[string]*Route
	orderedRouteNames []string
	currentRoute      *Route
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
	r.routes = make(map[string]*Route, 8)
	r.orderedRouteNames = make([]string, 0, 8)
	return r
}

// AddRoute adds an app.Route to this registry.
//
func (r *Registry) AddRoute(specs ...*Route) {

	/*
		cmds := make([]*commandSpec, 0, len(app.Route.Commands))
		for i, c := range app.Route.Commands {
			cmds[i] := &commandSpec {
			}
		}

		newroute := &Route{
			name:        route.Name,
			description: route.Help,
		}
	*/
	for _, spec := range specs {
		r.currentRoute = spec
		r.routes[spec.RouteName] = spec
		r.orderedRouteNames = append(r.orderedRouteNames, spec.RouteName)
	}
}

// Route specifies a new route to add to the registry.
func (r *Registry) Route(name, description string) *Registry {

	// Create the route spec.
	route := new(Route)
	route.RouteName = name
	route.Help = description
	route.Does = make([]*CommandSpec, 0, 4)

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
	spec := &CommandSpec{
		Name:    commandName,
		Command: cmd,
	}

	// Add command spec.
	r.currentRoute.Does = append(r.currentRoute.Does, spec)

	return r
}

// Using specifies a paramater to use for the most recently specified command
// as set by Does.
func (r *Registry) Using(name string) *Registry {
	// Look up the last command added.
	lastCommand := r.lastCommandAdded()

	// Create a new spec.
	spec := &ParamSpec{
		Name: name,
	}

	// Add it to the list.
	lastCommand.Using = append(lastCommand.Using, spec)
	return r
}

// WithDefault specifies the default value for the most recently specified
// parameter as set by Using.
func (r *Registry) WithDefault(value interface{}) *Registry {
	param := r.lastParamAdded()
	param.DefaultValue = value
	return r
}

// From sepcifies where to get the value from for the most recently specified
// paramater as set by Using.
func (r *Registry) From(fromVal ...string) *Registry {
	param := r.lastParamAdded()

	// This is sort of a hack. Really, we should make params.from a []string.
	param.From = strings.Join(fromVal, " ")
	return r
}

// Get the last parameter for the last command added.
func (r *Registry) lastParamAdded() *ParamSpec {
	cspec := r.lastCommandAdded()
	last := len(cspec.Using) - 1
	return cspec.Using[last]
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
	for _, cmd := range spec.Does {
		r.currentRoute.Does = append(r.currentRoute.Does, cmd)
	}
	return r
}

// RouteSpec gets a route cased on its name.
func (r *Registry) RouteSpec(routeName string) (spec *Route, ok bool) {
	spec, ok = r.routes[routeName]
	return
}

// Routes gets an unordered map of routes names to route specs.
//
// If order is important, use RouteNames to get the names (in order).
func (r *Registry) Routes() map[string]*Route {
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
func (r *Registry) lastCommandAdded() *CommandSpec {
	lastIndex := len(r.currentRoute.Does) - 1
	return r.currentRoute.Does[lastIndex]
}

// RouteDetails provides important information about a route.
type RouteDetails interface {
	Name() string
	Description() string
}

// Route describes a route.
//
// Formerly, this was a routeSpec.
type Route struct {
	RouteName, Help string
	Does            []*CommandSpec
}

func (r *Route) Name() string {
	return r.RouteName
}

func (r *Route) Description() string {
	return r.Help
}

// CommandSpec describes how a Route should treat a Command.
//
// This is used to link a command to a given route.
type CommandSpec struct {
	Name    string
	Command Command
	Using   Parameters
}

// ParamSpec is used to link a paramater (argument) to a Command.
//
// The router will pass in arguments to a command based on the
// ParamSpec.
type ParamSpec struct {
	Name         string
	DefaultValue interface{}
	From         string
}

type Parameters []*ParamSpec
type Commands []*CommandSpec

func Foo(foo ...CommandSpec) []*CommandSpec {
	res := make([]*CommandSpec, len(foo))
	for i, f := range foo {
		res[i] = &f
	}
	return res
}

// Param is a convenience function for building a new ParamSpec.
func Param(name string, defaultVal interface{}, from string) *ParamSpec {
	return &ParamSpec{
		Name:         name,
		DefaultValue: defaultVal,
		From:         from,
	}
}
