// Cookoo is a Chain-of-Command (CoCo) framework for writing
// applications.
//
// A chain of command framework works as follows:
//
// * A "route" is constructed as a chain of commands -- a series of
// single-purpose tasks that are run in sequence.
//
// * An application is composed of one or more routes.
//
// * Commands in a route communicate using a Context.
//
// * An application Router is used to receive a route name and then
// execute the appropriate chain of commands.
//
// To create a new Cookoo application, use cookoo.Cookoo(). This will
// configure and create a new registry, request router, and context.
// From there, use the Registry to build chains of commands, and then
// use the Router to execute chains of commands.
//
// Example:
// 
//    package main
//
//    import (
//      //This is the path to Cookoo
//      "github.com/masterminds/cookoo/src/cookoo"
//      "fmt"
//    )
//
//    func main() {
//      // Build a new Cookoo app.
//      registry, router, context := cookoo.Cookoo()
//
//      // Fill the registry.
//      registry.Route("TEST", "A test route").Does(HelloWorld, "hi") //...
//
//      // Execute the route.
//      router.HandleRequest("TEST", context, false)
//    }
//
//    func HelloWorld(cxt cookoo.Context, params *cookoo.Params) (interface{}, Interrupt) {
//      fmt.Println("Hello World")
//      return true, nil
//    }
//
// Unlike other CoCo implementations (like Pronto.js or Fortissimo),
// Cookoo commands are just functions.
//
// Interrupts:
//
// There are four types of interrupts that you may wish to return:
// - FatalError: This will stop the route immediately.
// - RecoverableError: This will allow the route to continue moving.
// - Stop: This will stop the current request, but not as an error.
// - Reroute: This will stop executing the current route, and switch to executing another route.
//
// To learn how to write Cookoo applications, you may wish to examine
// the small Skunk application: https://github.com/technosophos/skunk.
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
type Command func(cxt Context, params *Params) (interface{}, Interrupt)

// Generic return for a command.
// Generally, a command should return one of the following in the interrupt slot:
// - A FatalError, which will stop processing.
// - A RecoverableError, which will continue the chain.
// - A Reroute, which will cause a different route to be run.
type Interrupt interface {}

// A command can return a Reroute to tell the router to execute a different route.
type Reroute struct {
	route string
}
func (rr *Reroute) RouteTo() string {
	return rr.route
}

// Stop a route, but not as an error condition.
type Stop struct {}

// An error that should not cause the router to stop processing.
type RecoverableError struct {
	Message string
}
func (err *RecoverableError) Error() string {
	return err.Message
}

// A fatal error, which will stop the router from continuing a route.
type FatalError struct {
	Message string
}
func (err *FatalError) Error() string {
	return err.Message
}
