// Copyright 2013 Masterminds.

package cookoo;

import (
	"fmt"
)

type Registry struct {
	routes map[string]*routeSpec
	loggers []*loggerSpec
	currentRoute *routeSpec
	// datasources map[string]datasourceSpec
	// currentDS datasourceSpec
}

/*type Command struct {
	name string
}*/

// Execute a command and return a result.
type Command func(cxt *ExecutionContext, params map[string]*interface{}) interface{}

type Logger struct {
	impl interface{}
}

func (r *Registry) Init() *Registry {
	r.routes = make(map[string]*routeSpec, 8)
	return r
}

func (r *Registry) Route(name, description string) *Registry {

	// Create the route spec.
	route := new(routeSpec)
	route.name = name;
	route.description = description;
	route.commands = make([]*commandSpec, 0, 4)

	// Add the route spec.
	r.currentRoute = route
	r.routes[name] = route

	return r
}

func (r *Registry) Does(cmd Command, commandName string) *Registry {

	// Configure command spec.
	spec := new(commandSpec)
	spec.name = commandName
	spec.command = cmd

	// Add command spec.
	r.currentRoute.commands = append(r.currentRoute.commands, spec)

	return r
}

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

func (r *Registry) WithDefault(value *interface{}) *Registry {
	param := r.lastParamAdded()
	param.defaultValue = value
	return r
}

func (r *Registry) WithDefaultValue(value interface{}) *Registry {
	param := r.lastParamAdded()
	param.defaultValue = &value
	return r
}

func (r *Registry) From(fromVal string) *Registry {
	param := r.lastParamAdded()
	param.from = fromVal
	return r
}

// Get the last parameter for the last command added.
func (r *Registry) lastParamAdded() *paramSpec {
	cspec := r.lastCommandAdded()
	last := len(cspec.parameters) - 1
	return cspec.parameters[last]
}

func (r *Registry) Includes(route string) *Registry {
	fmt.Println("Need to finish for ", route)
	return r
}

// Add a logger to the registry.
// Once at least one logger has been added, the application can begin logging.
func (r *Registry) Logger(log *Logger, options map[string]string) *Registry {
	// Create a logger spec.
	spec := new(loggerSpec)
	spec.logger = log
	spec.options = options

	// Add the spec.
	r.loggers = append(r.loggers, spec)
	return r
}

func (r *Registry) Loggers() []*loggerSpec {
	return r.loggers
}

func (r *Registry) RouteSpec(routeName string) *routeSpec {
	return r.routes[routeName]
}

func (r *Registry) Routes() map[string]*routeSpec {
	return r.routes
}

// Look up the last command.
func (r *Registry) lastCommandAdded() *commandSpec {
	lastIndex := len(r.currentRoute.commands) - 1
	return r.currentRoute.commands[lastIndex]
}

type routeSpec struct {
	name, description string
	commands []*commandSpec
}

type commandSpec struct {
	name string
	command Command
	parameters []*paramSpec
}

type paramSpec struct {
	name string
	defaultValue *interface{}
	from string
}

type loggerSpec struct {
	logger *Logger
	options map[string]string
}
