package cli

import (
	"github.com/masterminds/cookoo"
	"testing"
)

func TestShowHelp(t *testing.T) {
	registry, router, context := cookoo.Cookoo();

	registry.Route("test", "Testing help.").Does(ShowHelp, "didShowHelp").
		Using("show").WithDefault(true).
		Using("summary").WithDefault("This is a summary")

		e := router.HandleRequest("test", context, false)

		if e != nil {
			t.Error("! Unexpected error.")
		}

		res := context.Get("didShowHelp").(bool)

		if !res {
			t.Error("! Expected help to be shown.")
		}
}
