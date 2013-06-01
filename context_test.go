// Copyright 2013 Masterminds

// This package provides the execution context for a Cookoo request.
package cookoo

import (
	"reflect"
	"runtime"
	"testing"
)

// An example datasource as can add to our store.
type ExampleDatasource struct {
	name string
}

// A simple equal function.
func equal(t *testing.T, a interface{}, b interface{}) {
	result := reflect.DeepEqual(a, b)
	if !result {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("Failed equals in %s:%d", file, line)
	}
}

func TestDatasource(t *testing.T) {
	foo := new(ExampleDatasource)
	foo.name = "bar"

	cxt := NewContext()

	cxt.AddDatasource("foo", foo)

	foo2 := cxt.Datasource("foo").(*ExampleDatasource)

	equal(t, foo, foo2)
	equal(t, "bar", foo2.name)

	cxt.RemoveDatasource("foo")

	equal(t, nil, cxt.Datasource("foo"))
}

func TestAddGet(t *testing.T) {
	cxt := NewContext()

	cxt.Add("test1", 42)
	cxt.Add("test2", "Geronimo!")
	cxt.Add("test3", func() string { return "Hello" })

	// Test Get
	equal(t, 42, cxt.Get("test1"))
	equal(t, "Geronimo!", cxt.Get("test2"))

	// Test has
	val, ok := cxt.Has("test1")
	if !ok {
		t.Error("! Failed to get 'test1'")
	}
	equal(t, 42, cxt.Get("test1"))

	_, ok = cxt.Has("test999")
	if ok {
		t.Error("! Unexpected result for 'test999'")
	}

	val, ok = cxt.Has("test3")
	fn := val.(func() string)
	if ok {
		equal(t, "Hello", fn())
	} else {
		t.Error("! Expected a function.")
	}

}

type LameStruct struct {
	stuff []string
}

func TestCopy(t *testing.T) {
	lame := new(LameStruct)
	lame.stuff = []string{"O", "Hai"}
	c := NewContext()
	c.Add("a", lame)
	c.Add("b", "This is the song that never ends")

	c2 := c.Copy()

	c.Add("c", 1234)

	if c.Len() != 3 {
		t.Error("! Canary failed. c should be 3")
	}

	if c2.Len() != 2 {
		t.Error("! c2 should be 2.")
	}

	c.Add("b", "FOO")
	if c2.Get("b") == "FOO" {
		t.Error("! b should not have changed in C2.")
	}

	lame.stuff[1] = "Noes"

	v1 := c2.Get("a").(*LameStruct)
	if v1.stuff[1] != "Noes" {
		t.Error("! Expected shallow copy of array. Got ", v1)
	}
}
