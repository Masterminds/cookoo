// Copyright 2013 Masterminds

// This package provides the execution context for a Cookoo request.
package context

import (
	"github.com/bmizerany/assert"
	"testing"
)

// An example datasource as can add to our store.
type ExampleDatasource struct {
	name string
}

func TestDatasource(t *testing.T) {
	foo := new(ExampleDatasource)
	foo.name = "bar"

	cxt := NewExecutionContext()

	cxt.AddDatasource("foo", foo)

	foo2 := cxt.Datasource("foo").(*ExampleDatasource)

	assert.Equal(t, foo, foo2)
	assert.Equal(t, "bar", foo2.name)
}
