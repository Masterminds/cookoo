// Copyright 2013 Masterminds

// This package provides the execution context for a Cookoo request.
package cookoo

import (
	"github.com/bmizerany/assert"
	"testing"
	//"fmt"
	//"reflect"
)

// An example datasource as can add to our store.
type ExampleDatasource struct {
	name string
}

func TestDatasource(t *testing.T) {
	foo := new(ExampleDatasource)
	foo.name = "bar"

	cxt := NewContext()

	cxt.AddDatasource("foo", foo)

	foo2 := cxt.Datasource("foo").(*ExampleDatasource)

	assert.Equal(t, foo, foo2)
	assert.Equal(t, "bar", foo2.name)

	cxt.RemoveDatasource("foo")

	assert.Equal(t, nil, cxt.Datasource("foo"))
}

func TestAddGet(t *testing.T) {
	cxt := NewContext()

	cxt.Add("test1", 42)
	cxt.Add("test2", "Geronimo!")
	cxt.Add("test3", func() string { return "Hello" })

	// Test Get
	assert.Equal(t, 42, cxt.Get("test1"))
	assert.Equal(t, "Geronimo!", cxt.Get("test2"))

	// Test has
	val, ok := cxt.Has("test1")
	if !ok {
		t.Error("! Failed to get 'test1'")
	}
	assert.Equal(t, 42, cxt.Get("test1"))

	_, ok = cxt.Has("test999")
	if ok {
		t.Error("! Unexpected result for 'test999'")
	}

	val, ok = cxt.Has("test3")
	fn := val.(func() string)
	if ok {
		assert.Equal(t, "Hello", fn())
	} else {
		t.Error("! Expected a function.")
	}

}
