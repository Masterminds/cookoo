package cli

import (
	"github.com/masterminds/cookoo"
	"testing"
	"bytes"
	"strings"
	"flag"
	//"fmt"
)

func Barf(cxt cookoo.Context, params *cookoo.Params) (interface{}, cookoo.Interrupt) {
	return nil, cookoo.FatalError{"Intentional fail!"}
}

func TestShowHelp(t *testing.T) {
	registry, router, context := cookoo.Cookoo();

	var out bytes.Buffer

	registry.Route("test", "Testing help.").Does(ShowHelp, "didShowHelp").
		Using("show").WithDefault(true).
		Using("writer").WithDefault(&out).
		Using("summary").WithDefault("This is a summary.").
		Does(Barf, "Fail if help doesn't stop.")


	e := router.HandleRequest("test", context, false)

	if e != nil {
		t.Error("! Unexpected error.")
	}

	res := context.Get("didShowHelp").(bool)

	if !res {
		t.Error("! Expected help to be shown.")
	}

	msg := out.String()
	if !strings.Contains(msg, "SUMMARY\n") {
		t.Error("! Expected 'summary' as a header.")
	}
	if !strings.Contains(msg, "This is a summary.") {
		t.Error("! Expected 'This is a summary' to be in the output. Got ", msg)
	}
}

func TestParseArgs(t *testing.T) {
	registry, router, cxt := cookoo.Cookoo();

	flags := flag.NewFlagSet("test flags", flag.ContinueOnError)
	flags.String("foo", "binky", "Test foo flag.")
	flags.Bool("baz", false, "Baz flag")
	flags.Int("unused", 123, "Unused int flag.")

	registry.Route("test", "Testing parse arguments.").
		Does(ParseArgs, "args").
		Using("args").WithDefault([]string{"-foo", "bar", "-baz", "arg1"}).
		Using("flagset").WithDefault(flags)

	if router.HandleRequest("test", cxt, false) != nil {
		t.Error("! Request failed.")
		return
	}

	foo := cxt.Get("foo").(*flag.Flag).Value.String()
	if foo != "bar" {
		t.Error("Expected 'bar'; got ", foo)
		return
	}

	bazO, ok := cxt.Has("baz")
	if !ok {
		t.Error("Expected to find 'baz' in context.")
		return
	}
	baz := bazO.(*flag.Flag).Value
	// fmt.Printf("baz is %v", baz)
	if baz.String() != "true" {
		t.Error("Expected 'baz' to be true. Got false.")
		return
	}

	unused := cxt.Get("unused").(*flag.Flag).Value
	if unused.String() != "123" {
		t.Error("Expected 'unused' to be int 123. Got ", unused)
		return
	}
}
