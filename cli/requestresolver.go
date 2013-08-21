package cli

import (
	"flag"
	"github.com/masterminds/cookoo"
)

type RequestResolver struct {
	registry *cookoo.Registry
}

func (r *RequestResolver) Init(registry *cookoo.Registry) {
	r.registry = registry
}

func (r *RequestResolver) Resolve(path string, cxt cookoo.Context) string {
	// Parse out any flags. Maybe flag specs are in context?

	flagset, ok := cxt.Has("globalFlags")
	if ok {
		r.addFlagsToContext(flagset.(*flag.FlagSet), cxt)
	}

	// Parse argv[0] as subcommand
	// Adjust the rest of argv to pass into the coco?

	// If all else fails, just return the unlatered path.
	return path
}

func (r *RequestResolver) addFlagsToContext(flagset *flag.FlagSet, cxt cookoo.Context) {
	store := func(f *flag.Flag) {
		cxt.Add(f.Name, f)
	}

	flagset.VisitAll(store)
}
