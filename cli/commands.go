package cli

import (
	"github.com/masterminds/cookoo"
	"fmt"
	"os"
	"io"
	"strings"
	"flag"
)

// Parse arguments for a "subcommand"
// 
// The cookoo.cli.RequestResolver allows you to specify global level flags. This command
// allows you to augment those with subcommand flags. Example:
//
// 		$ myprog -foo=yes subcommand -bar=no
//
// In the above example, `-foo` is a global flag (set before the subcommand), while
// `-bar` is a local flag. It is specific to `subcommand`. This command lets you parse
// an arguments list given a pointer to a `flag.FlagSet`.
//
// Like the cookoo.cli.RequestResolver, this will place the parsed params directly into the
// context. For this reason, you ought not use the same flag names at both global and local
// flag levels. (The local will overwrite the global.)
//
// Params:
// - args: (required) A slice of arguments. Typically, this is `cxt:args` as set by
// 		cookoo.cli.RequestResolver.
// - flagset: (required) A set if flags (see flag.FlagSet) to parse.
//
// A slice of all non-flag arguments remaining after the parse are returned into the context.
//
// For example, if ['-foo', 'bar', 'some', 'other', 'data'] is passed in, '-foo' and 'bar' will
// be parsed out, while ['some', 'other', 'data'] will be returned into the context. (Assuming, of
// course, that the flag definition for -foo exists, and is a type that accepts a value).
//
// Thus, you will have `cxt:foo` available (with value `bar`) and everything else will be available
// in the slice under this command's context entry.
func ParseArgs(cxt cookoo.Context, params *cookoo.Params) (interface{}, cookoo.Interrupt) {
	params.Requires("args", "flagset")
	flagset := params.Get("flagset", nil).(*flag.FlagSet)
	args := params.Get("args", nil).([]string)

	flagset.Parse(args)
	addFlagsToContext(flagset, cxt)
	return flagset.Args(), nil

}

// Show help.
// This command is useful for placing at the front of a CLI "subcommand" to have it output
// help information. It will only trigger when "show" is set to true, so another command
// can, for example, check for a "-h" or "-help" flag and set "show" based on that.
//
// Params:
// - show (bool): If `true`, show help.
// - summary (string): A one-line summary of the command.
// - description (string): A short description of what the command does.
// - usage (string): usage information.
// - writer (Writer): The location that this will write to. Default is os.Stdout
func ShowHelp(cxt cookoo.Context, params *cookoo.Params) (interface{}, cookoo.Interrupt) {
	showHelp := params.Get("show", false).(bool)
	writer := params.Get("writer", os.Stdout).(io.Writer)

	if showHelp {
		displayHelp([]string{"summary", "description", "usage"}, params, writer)
		return true, new(cookoo.Stop)
	}

	return false, nil
}

func displayHelp(keys []string, params *cookoo.Params, out io.Writer) {
	for i := range keys {
		key := keys[i]
		msg, ok := params.Has(key)
		if ok {
			fmt.Fprintf(out, "%s\n\n%s\n", strings.ToUpper(key), msg)
		}
	}
}
