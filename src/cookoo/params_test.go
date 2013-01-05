package cookoo

import (
	"testing"
)

func TestParams(t *testing.T) {
	params := Params{
		"Test": 123,
		"Test2": "Hello",
		"Test3": NewContext(),
	}

	if v, ok := params.Has("Test"); !ok {
		t.Error("Expected to find 123, got NADA")
	} else if v != 123 {
		t.Error("! Expected 123, got ", v)
	}

	// A really lame validator.
	fn := func(value interface{}) bool {
		return true
	}

	// Test the validator.
	if v, ok := params.Validate("Test2", fn); !ok {
		t.Error("! Expected a valid string.")
	} else if v != "Hello" {
		t.Error("! Expected 'Hello', got ", v)
	}

	alwaysFails := func(value interface{}) bool {
		return false
	}
	// Test the validator.
	if _, ok := params.Validate("Test2", alwaysFails); ok {
		t.Error("! Expected a failed validation.")
	}
}
