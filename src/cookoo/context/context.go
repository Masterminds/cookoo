// Copyright 2013 Masterminds

// This package provides the execution context for a Cookoo request.
package context

// An empty interface defining a context value.
// Semantically, this is the same as interface{}
type ContextValue interface{}

type ExecutionContext struct {
	datasources map[string]interface{} // Datasources are things like MySQL connections.

	// The Context values.
	values map[string]interface{}
}

func NewExecutionContext() *ExecutionContext {
	cxt := new(ExecutionContext).Init()
	return cxt
}

func (cxt *ExecutionContext) Init() *ExecutionContext {
	cxt.datasources = make(map[string]interface{})
	cxt.values = make(map[string]interface{})
	return cxt
}

// Add a name/value pair to the context.
func (cxt *ExecutionContext) Add(name string, value ContextValue) {
	cxt.values[name] = value
}

// Given a name, return the corresponding value from the context.
func (cxt *ExecutionContext) Get(name string) ContextValue {
	return cxt.values[name]
}

// A special form of Get that also returns a flag indicating if the value is found.
// This fetches the value and also returns a flag indicating if the value was
// found. This is useful in cases where the value may legitimately be 0.
func (cxt *ExecutionContext) Has(name string) (value ContextValue, found bool) {
	value, found = cxt.values[name]
	return;
}

// Get a datasource from the map of datasources.
// A datasource (e.g., a connection to a database) is retrieved as an interface
// so its type will need to be specified before it can be used. Take an example
// of the variable foo that is a struct of type Foo.
// foo = cxt.Datasource("foo").(*Foo)
func (cxt *ExecutionContext) Datasource(name string) interface{} {
	return cxt.datasources[name]
}

// Add a datasource to the map of datasources.
// A datasource is typically something like a connection to a database that you
// want to keep open persistently and share between requests. To add a datasource
// to the map just add it with a name. e.g. cxt.AddDatasource("mysql", foo) where
// foo is the struct for the datasource.
func (cxt *ExecutionContext) AddDatasource(name string, ds interface{}) {
	cxt.datasources[name] = ds
}
