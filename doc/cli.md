# Writing a CLI

Here is a basic scaffold of a CLI application that uses subcommands.
Here, we're building the command `foo` that has the subcommand `bar`,
which can be invoked like this from tthe commandline: `foo bar`.

foo.go:
```go
package main

import(
  "github.com/masterminds/cookoo"
  "github.com/masterminds/cookoo/cli"
  "fmt"
  "os"
)

func main() {
	// Start a cookoo app.
	reg, router, cxt := cookoo.Cookoo();

	// Put the arguments into the context.
	cxt.Add("os.Args", os.Args)

	// Create help text
	reg.Route("help", "Show application-scoped help.").
		Does(cli.ShowHelp, "help").
			Using("show").WithDefault(true).
			Using("summary").WithDefault("This is the help text.")

	// Handle the "bar" subcommand
	reg.Route("bar", "Do something").
		Does(MyBarCommand, "bar")

	// This is the main runner. It proxies to subcommands.
	reg.Route("run", "Run the app.").
		Does(cli.RunSubcommand, "sub").
			Using("args").From("cxt:os.Args").
			Using("default").WithDefault("help").
			Using("ignoreRoutes").WithDefault([]string{"run"})

	// This starts the app.	If a fatal error occurs, we
	// display the error.
	e := router.HandleRequest("run", cxt, true)
	if e != nil {
		fmt.Printf("Error: %s\n", e)
	}
}

// This is our command.
func MyBarCommand(cxt cookoo.Context, params *cookoo.Params) (interface{}, cookoo.Interrupt) {
	fmt.Printf("OH HAI")
	return nil, nil
}
```

When we run `foo bar` (or `go run foo.go bar`), here's what happens:

1. main() is run. It creates a Cookoo app, defines the registry, and
   then runs `router.HandleRequest("run", ...)`
2. The "run" route is run. This reads the `os.Args` and sees that it
   should execute the subcommand `bar`.
3. The "bar" route is run, which executes it's one command:
   `MyBarCommand`.
4. `MyBarCommand` runs, printing "OH HAI" to stdout.

If you were to run `foo` or `foo help`, then the chain would execute
like this:

1. main() runs, and passes to `router.HandleRequest("run"...)
2. "run" will execute `RunSubcommand`, which will resolve the subcommand
   to "help" (which is the default target).
3. The "help" route will be run, which will print out simple help:

```
go run foo.go
SUMMARY

This is the help text.
```

## Where to from here?

From this starting point, you should be able to assemble your own
routes, where each new route represents a subcommand.
