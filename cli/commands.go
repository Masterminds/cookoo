package cli

import (
	"github.com/masterminds/cookoo"
	"fmt"
)

func ShowHelp(cxt cookoo.Context, params *cookoo.Params) (interface{}, cookoo.Interrupt) {
	showHelp := params.Get("show", false).(bool)

	if !showHelp {
		displayHelp([]string{"summary", "description", "usage"}, params)
		return false, nil
	}

	return true, new(cookoo.Stop)
}

func displayHelp(keys []string, params *cookoo.Params) {
	for i := range keys {
		key := keys[i]
		msg, ok := params.Has(key)
		if ok {
			fmt.Printf("%s\n\n%s\n", key, msg)
		}
	}
}
