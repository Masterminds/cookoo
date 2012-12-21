// Copyright 2013 Masterminds

// This package provides the execution context for a Cookoo request.
package context

type Datasource struct {
}

type ExecutionContext struct {
	// Need the following:
  // Datasources -- probably a hashtable
  // Context vars -- hashtable
}

func (cxt *ExecutionContext) Init() *ExecutionContext {
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
