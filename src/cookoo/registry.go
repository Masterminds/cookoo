// Copyright 2013 Masterminds.

package cookoo;

import (
	"fmt"
)

type Registry struct {
	var routes map[string]routeSpec
	var loggers []loggerSpec
	var currentRoute routeSpec
	// var datasources map[string]datasourceSpec
	// var currentDS datasourceSpec
}

type Command struct {
	var name string
}

type Logger struct {
	var impl interface{}
}

func (r *Registry) Route(name, description string) *Registry {

	// Create the route spec.
	route := new(routeSpec)
	route.name = name;
	route.description = description;
	route.commands = make(commandSpec, 2, 8)

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
	append(r.currentRoute.commands, spec)

	return r
}

func (r *Registry) Using(name string) *Registry {
	// Look up the last command added.
	lastCommand := r.lastCommandAdded()

	// Create a new spec.
	spec := new(paramSpec)
	spec.name = name

	// Add it to the list.
	append(lastCommand.paramaters, spec)
	return r
}

func (r *Registry) WithDefault(value interface{}) *Registry {
	param := r.lookupLastParam()
	param.defaultValue = value
	return r
}

func (r *Registry) From(fromVal string) *Registry {
	param := r.lookupLastParam()
	param.from = fromVal
	return r
}

// Get the last parameter for the last command added.
func (r *Registry) lookupLastParam() *paramSpec {
	cspec := r.lookupLastCommand()
	last := len(cspec.parameters) - 1
	return cspec.parameters[last]
}

func (r *Registry) Includes(route string) *Registry {
	fmt.println("Need to finish for ", route)
	return r
}

// Add a logger to the registry.
// Once at least one logger has been added, the application can begin logging.
func (r *Registry) Logger(log Logger, options map[string]string) *Registry {
	// Create a logger spec.
	spec := new(loggerSpec)
	spec.logger = log
	spec.options = options

	// Add the spec.
	append(r.loggers, spec)
	return r
}

func (r *Registry) Loggers() []Logger {
	return r.loggers
}

func (r *Registry) RouteSpec(routeName string) routeSpec {
	return r.routeSpec[routeName]
}

func (r *Registry) Routes() map[string]routeSpec {
	return r.routes
}

// Look up the last command.
func (r *Registry) lastCommandAdded() *commandSpec {
	lastIndex := len(r.currentRoute.commands) - 1
	return r.currentRoute.commands[lastIndex]
}

type routeSpec struct {
	var name, description string
	var commands []commandSpec
}

type commandSpec struct {
	var name string
	var command Command
	var parameters []paramSpec
}

type paramSpec struct {
	var name string
	var defaultValue interface{}
	var from string
}

type loggerSpec struct {
	var logger Logger
	var options map[string]string
}
