package cli

import (
	"flag"
	"github.com/masterminds/cookoo"
	//"os"
	"strings"
	//"fmt"
)

type RequestResolver struct {
	registry *cookoo.Registry
}

func (r *RequestResolver) Init(registry *cookoo.Registry) {
	r.registry = registry
}

func (r *RequestResolver) Resolve(path string, cxt cookoo.Context) (string, error) {
	// Parse out any flags. Maybe flag specs are in context?

	flagsetO, ok := cxt.Has("globalFlags")
	if !ok {
		// No args to parse. Just return path.
		return path, nil
	}
	flagset := flagsetO.(*flag.FlagSet)
	flagset.Parse(strings.Split(path, " "));
	addFlagsToContext(flagset, cxt)
	args := flagset.Args()

	// This is a failure condition... Need to fix Cookoo to support error return.
	if len(args) == 0 {
		return path, &cookoo.RouteError{"Could not resolve route " + path}
	}

	// Add the rest of the args to the context.
	cxt.Add("args", args[1:])

	// Parse argv[0] as subcommand
	return args[0], nil
}

func addFlagsToContext(flagset *flag.FlagSet, cxt cookoo.Context) {
	store := func(f *flag.Flag) {
		// fmt.Printf("Storing %s in context with value %s.\n", f.Name, f.Value.String())
		cxt.Add(f.Name, f)
	}

	flagset.VisitAll(store)
}
