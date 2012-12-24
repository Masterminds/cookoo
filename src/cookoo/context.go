// Copyright 2013 Masterminds

// This package provides the execution context for a Cookoo request.
package context

// TODO: Turn Datasource into an interface?
type Datasource struct {
}

type ExecutionContext struct {
	datasources map[string]Datasource // Datasources are things like MySQL connections.
	// Need the following:
	// Context vars -- hashtable
}

func NewExecutionContext() *ExecutionContext {
	cxt := new(ExecutionContext).Init()
	return cxt
}

func (cxt *ExecutionContext) Init() *ExecutionContext {
	cxt.datasources = make(map[string]Datasource)
	return cxt
}

func (cxt *ExecutionContext) Add(name string, value string) {
}

func (cxt *ExecutionContext) Get(name string) {
}

func (cxt *ExecutionContext) Datasource() *Datasource {
	return null
}

func (cxt *ExecutionContext) AddDatasource(*Datasource) {
}
