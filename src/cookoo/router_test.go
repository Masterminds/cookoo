package cookoo

import (
	"testing"
)

type FakeRequestResolver struct {
	BasicRequestResolver
}

func TestSetResolver (t *testing.T) {
	registry := new(Registry)
	r := new(Router)
	r.Init(registry)

	//var resolver *RequestResolver = r.RequestResolver()

	resolver := new(FakeRequestResolver)
	r.SetRequestResolver(resolver)
	resolver, ok := r.RequestResolver().(*FakeRequestResolver)

	if !ok {
		t.Error("! Resolver is not a FakeRequestResolver.")
	}
}
