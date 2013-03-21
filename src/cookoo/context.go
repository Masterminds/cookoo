// Copyright 2013 Masterminds

// This package provides the execution context for a Cookoo request.
package cookoo

// Describes a context.
type Context interface {
	// Add a name/value pair to the context.
	Add(string, ContextValue)
	// Given a name, get a value from the context.
	Get(string) ContextValue
	// Given a name, check if the key exists, and if it does return the value.
	Has(string) (ContextValue, bool)
	// Get a datasource by name.
	Datasource(string) Datasource
	// Get a map of all datasources.
	Datasources() map[string]Datasource
	// Check if a datasource exists, and return it if it does.
	HasDatasource(string) (Datasource, bool)
	// Add a datasource.
	AddDatasource(string, Datasource)
	// Remove a datasource from the context.
	RemoveDatasource(string)
	// Get the length of the context. This is the number of context values.
	// Datsources are not counted.
	Len() int
	// Make a shallow copy of the context.
	Copy() Context
	// Get the content (no datasources) as a map.
	AsMap() map[string]ContextValue
}

// An empty interface defining a context value.
// Semantically, this is the same as interface{}
type ContextValue interface{}

// An empty interface defining a Datasource.
// Semantically, this is the same as interface{}
type Datasource interface{}

type ExecutionContext struct {
	datasources map[string]Datasource // Datasources are things like MySQL connections.

	// The Context values.
	values map[string]ContextValue
}

// A datasource that can retrieve values by (string) keys.
// Datsources can be just about anything. But a key/value datasource
// can be used for a special purpose. They can be accessed in From()
// clauses in a registry configuration.
type KeyValueDatasource interface {
	Value(key string) interface{}
}

func NewContext() Context {
	cxt := new(ExecutionContext).Init()
	return cxt
}

func (cxt *ExecutionContext) Init() *ExecutionContext {
	cxt.datasources = make(map[string]Datasource)
	cxt.values = make(map[string]ContextValue)
	return cxt
}

// Add a name/value pair to the context.
func (cxt *ExecutionContext) Add(name string, value ContextValue) {
	cxt.values[name] = value
}

func (cxt *ExecutionContext) AsMap() map[string]ContextValue {
	return cxt.values
}

// Given a name, return the corresponding value from the context.
func (cxt *ExecutionContext) Get(name string) ContextValue {
	return cxt.values[name]
}

// Get a map of all name/value pairs in the present context.
func (cxt *ExecutionContext) GetAll() map[string]ContextValue {
	return cxt.values
}

// A special form of Get that also returns a flag indicating if the value is found.
// This fetches the value and also returns a flag indicating if the value was
// found. This is useful in cases where the value may legitimately be 0.
func (cxt *ExecutionContext) Has(name string) (value ContextValue, found bool) {
	value, found = cxt.values[name]
	return
}

// Get a datasource from the map of datasources.
// A datasource (e.g., a connection to a database) is retrieved as an interface
// so its type will need to be specified before it can be used. Take an example
// of the variable foo that is a struct of type Foo.
// foo = cxt.Datasource("foo").(*Foo)
func (cxt *ExecutionContext) Datasource(name string) Datasource {
	return cxt.datasources[name]
}

func (cxt *ExecutionContext) Datasources() map[string]Datasource {
	return cxt.datasources
}

// Check whether the named datasource exists, and return it if it does.
func (cxt *ExecutionContext) HasDatasource(name string) (Datasource, bool) {
	value, found := cxt.datasources[name];
	return value, found;
}

// Add a datasource to the map of datasources.
// A datasource is typically something like a connection to a database that you
// want to keep open persistently and share between requests. To add a datasource
// to the map just add it with a name. e.g. cxt.AddDatasource("mysql", foo) where
// foo is the struct for the datasource.
func (cxt *ExecutionContext) AddDatasource(name string, ds Datasource) {
	cxt.datasources[name] = ds
}

func (cxt *ExecutionContext) RemoveDatasource(name string) {
	delete(cxt.datasources, name)
}

func (cxt *ExecutionContext) Len() int {
	return len(cxt.values)
}

// Copy the context into a new context.
func (cxt *ExecutionContext) Copy() Context {
	newCxt := NewContext()
	vals := cxt.GetAll()
	ds := cxt.Datasources()

	for k, v := range vals {
		newCxt.Add(k, v)
	}

	for k, datasource := range ds {
		newCxt.AddDatasource(k, datasource)
	}

	return newCxt;
}
