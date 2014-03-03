package cookoo

import (
	"testing"
)

func TestCookooForCoCo(t *testing.T) {
	registry, router, cxt := Cookoo()

	cxt.Put("Answer", 42)

	lifeUniverseEverything := cxt.Get("Answer", nil)

	if lifeUniverseEverything != 42 {
		t.Error("! Context is not working.")
	}

	registry.Route("foo", "test")

	ok := router.HasRoute("foo")

	if !ok {
		t.Error("! Router does not have 'foo' route.")
	}
}
