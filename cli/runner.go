package cli

import (
	"github.com/Masterminds/cookoo"
	"flag"
	"fmt"
	"os"
)


// New creates a new CLI instance.
//
// It takes a flagset for parsing command line options, and creates a new
// CLI application initialized. Flags are placed into the context.
//
// By default, the `cookoo.BasicRequestResolver` is used for resolving request.
// This works well with the subcommand model of delegating commands.
//
// If a '@startup' route is inserted into the registry, it will be run first
// upon any call to `Router.HandleRequest`. If not, the default startup will
// be run. That routine includes displaying help text if the -h or -help
// flags are passed in.
//
// The `summary` is a one line explanation of the program used in help text.
//
// The `usage` is a detailed help message, often several paragraphs.
//
// The `globalFlags` are a `flag.FlagSet` for the top level of this program. If
// you use subcommands and want subcommand-specific flags, use the `ParseArgs`
// command in this package.
//
// Typical usage:
//
// 	package main
//
// 	import(
// 		"github.com/Masterminds/cookoo"
// 		"github.com/Masterminds/cookoo/cli"
// 		"flag"
// 	)
//
// 	var Summary := "A program that does stuff."
// 	var Description := `Full text description goes here.`
// 	func main() {
// 		flags := flag.FlagSet("global", flag.PanicOnError)
// 		// Define any flags here...
// 		flags.Bool("h", false, "Show help text")
//
// 		reg, router, cxt := cookoo.Cookoo()
// 		reg.Route("hello", "Does nothing")
//
// 		cli.New(reg, router, cxt).Help(Summary, Description, flags).Run("hello")
// 	}
//
// The simple program above can be run any of these ways:
//
// 	$ mycli       # Will run 'hello', which does nothing
//  $ mycli -h    # Will show help
//
// If we were to substitute `RunSubcommand` instead of `Run`:
//
// 	cli.New(reg, router, cxt).Help(Summary, Description, flags).RunSubcommand()
//
// The above would do the following:
// 	$ mycli       # Will show help
//  $ mycli -h    # Will show help
//  $ mycli help  # Will show help
//  $ mycli hello # Will run hte "hello" route, which does nothing.
//
func New(reg *cookoo.Registry, router *cookoo.Router, cxt cookoo.Context) *Runner {
	return &Runner{reg: reg, router: router, cxt: cxt}
}

// Runner is a CLI runner.
//
// It provides a nice abstraction for simply and easily running CLI
// commands.
type Runner struct {
	reg *cookoo.Registry
	router *cookoo.Router
	cxt cookoo.Context

	summary, usage string
	flags *flag.FlagSet
}

// Help sets the help text and support flags for the app.
//
// It is strongly advised to use this function for all CLI runner apps.
func (r *Runner) Help(summary, usage string, flags *flag.FlagSet) *Runner {
	r.summary = summary
	r.usage = usage
	r.flags = flags
	return r
}

func (r *Runner) startup() {

	if r.flags == nil {
		r.flags = flag.NewFlagSet("globalFlags", flag.PanicOnError)
		r.flags.Bool("h", false, "Show this help text.")
		r.flags.Bool("help", false, "Show this help text.")
	}

	r.cxt.Put("globalFlags", r.flags)
	r.cxt.Put("os.Args", os.Args)

	// Allow route to be overwritten.
	if _, ok := r.reg.RouteSpec("@startup"); !ok {
		r.reg.Route("@startup", "Prepare to run a route.").
			Does(ShiftArgs, "_").Using("n").WithDefault(1).
			Does(ParseArgs, "runner.Args").
				Using("flagset").WithDefault(r.flags).
				Using("args").From("cxt:os.Args").
			Does(ShowHelp, "help").
				Using("show").From("cxt:h").
				Using("summary").WithDefault(r.summary).
				Using("usage").WithDefault(r.usage).
				Using("flags").WithDefault(r.flags).
			Does(ShowHelp, "-help"). // Stupid hack. FIXME.
				Using("show").From("cxt:help").
				Using("summary").WithDefault(r.summary).
				Using("usage").WithDefault(r.usage).
				Using("flags").WithDefault(r.flags)
	}
	if _, ok := r.reg.RouteSpec("@subcommand"); !ok {
		r.reg.Route("@subcommand", "Startup and run subcommand").
			Includes("@startup").
			Does(RunSubcommand, "subcommand").
			Using("default").WithDefault("help").
			Using("offset").WithDefault(0).
			Using("args").From("cxt:runner.Args").
			Using("ignoreRoutes").WithDefault([]string{"@startup", "@subcommand"})
	}

	if _, ok := r.reg.RouteSpec("help"); !ok {
		r.reg.Route("help", "Basic help command.").
			Does(ShowHelp, "help").
			Using("show").WithDefault(true).
			Using("summary").WithDefault(r.summary).
			Using("usage").WithDefault(r.usage).
			Using("flags").WithDefault(r.flags)
	}
}

// Run runs a given route.
//
// It first runs the '@startup' route, and then runs whatever the named route
// is.
//
// If the flags `-h` or `-help` are specified, then the presence of those
// flags will automatically trigger help text.
//
// Additionally, the command `help` is predefined to generate help text.
func (r *Runner) Run(route string) error {
	r.startup()
	if err := r.router.HandleRequest("@startup", r.cxt, false); err != nil {
		fmt.Printf("Failed to startup: %s", err)
		os.Exit(1)
		//return err
	}
	if r.cxt.Get("help", false).(bool) {
		return nil
	}
	// FIXME: Hack
	if r.cxt.Get("-help", false).(bool) {
		return nil
	}
	return r.router.HandleRequest(route, r.cxt, true)
}

// RunSubcommand uses the first non-flag argument in args as a route name.
//
// For example:
// 	$ mycli -foo=bar myroute
//
// The above will see 'myroute' as a subcommand, and match it to a route named
// 'subcommand'. In the event that the route is not present, the help text
// will be displayed.
//
//
func (r *Runner) RunSubcommand() error {
	r.startup()
	return r.router.HandleRequest("@subcommand", r.cxt, false)
}
