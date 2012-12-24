// Copyright 2013 Masterminds

// This package provides the execution context for a Cookoo request.
package context

type ExecutionContext struct {
	datasources map[string]interface{} // Datasources are things like MySQL connections.
	// Need the following:
	// Context vars -- hashtable
}

func NewExecutionContext() *ExecutionContext {
	cxt := new(ExecutionContext).Init()
	return cxt
}

func (cxt *ExecutionContext) Init() *ExecutionContext {
	cxt.datasources = make(map[string]interface{})
	return cxt
}

func (cxt *ExecutionContext) Add(name string, value string) {
}

func (cxt *ExecutionContext) Get(name string) {
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
