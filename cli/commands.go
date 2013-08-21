package cli

import (
	"github.com/masterminds/cookoo"
	"fmt"
	"os"
	"io"
	"strings"
)

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
