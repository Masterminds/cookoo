// Copyright 2013 Masterminds

package cookoo

import (
	cio "github.com/masterminds/cookoo/io"
	"io"
	"log"
)

// A Context is a collection of data that is associated with the current
// request.
//
// Contexts are used to exchange information from command to command inside
// of a particular chain of commands (a route). Commands may access the
// data inside of a context, and may also modify a context.
//
// A context maintains two different types of data: *context variables* and
// *datasources*.
//
// Context variables are data that can be passed, in current form, from
// command to command -- analogous to passing variables via parameters in
// function calls.
//
// Datasources are (as the name implies) sources of data. For example, a
// database, a file, a cache, and a key-value store are all datasources.
//
// For long-running apps, it is generally assumed (though by no means
// required) that datasources are "long lived" and context variables are
// "short lived." While modifying a data source may impact other requests,
// generally it is safe to assume that modifying a variable is localized to
// the particular request.
//
// Correct Usage: A Word of Warning
// =================
//
// The Cookoo system was designed around the theory that commands should
// generally work with datasources *directly* and context variables
// *indirectly*. Context variables should generally be passed into a command
// via a cookoo.Param. And a command generally should return a value that
// can then be placed into the context on its behalf.
//
// The reason for this design is that it then makes it much easier for higher-
// level programming, such as changing input or modifying output at the
// registry level, not within the commands themselves.
//
// Datasources, on the other hand, are designed to be leveraged primarily by
// commands. This involves a layer of conventionality, but it also pushes
// data access logic into the commands where it belongs.
//
// So, for example, a SQL-based datasource should be *declared* at the top
// level of a program (where it will be added to the context), but the actual
// interaction with that datasource should happen inside of commands themselves,
// not at the registry level.
type Context interface {
	// Add a name/value pair to the context.
	Add(string, ContextValue)
	// Given a name, get a value from the context.
	//
	// Get requires a default value (which may be nil).
	//
	// Example:
	// 	ip := cxt.Get("ip", "127.0.0.1").(string)
	//
	// Contrast this usage with that of cxt.Has(), which may be used for more
	// traditional field checking:
	//
	// Example:
	// 	ip, ok := cxt.Has("ip")
	// 	if !ok {
	// 		// do something error-ish
	// 	}
	// 	ipStr := ip.(string)
	//
	// The cxt.Get() call avoids the cumbersome check/type-assertion combo
	// that occurs with cxt.Has().
	Get(string, interface{}) ContextValue
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
	// Get a logger.
	Logger(name string) (io.Writer, bool)
	// Add a logger.
	AddLogger(name string, logger io.Writer)
	// Remove a logger.
	RemoveLogger(name string)
	// Send a log with a prefix.
	Log(prefix string, v ...interface{})
	// Send a log and formatting string with a prefix.
	Logf(prefix string, format string, v ...interface{})
}

// An empty interface defining a context value.
// Semantically, this is the same as interface{}
type ContextValue interface{}

// An empty interface defining a Datasource.
// Semantically, this is the same as interface{}
type Datasource interface{}

// The core implementation of a Context.
//
// An ExecutionContext is an unordered map-based context.
type ExecutionContext struct {
	datasources map[string]Datasource // Datasources are things like MySQL connections.

	// The Context values.
	values map[string]ContextValue

	loggers          io.Writer
	loggerRegistered bool
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
	cxt.loggers = cio.NewMultiWriter()
	cxt.loggerRegistered = false
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
func (cxt *ExecutionContext) Get(name string, defaultValue interface{}) ContextValue {
	val, ok := cxt.values[name]
	if !ok {
		return defaultValue
	}
	return val
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
	value, found := cxt.datasources[name]
	return value, found
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

func (cxt *ExecutionContext) Logger(name string) (io.Writer, bool) {
	writer, found := cxt.loggers.(*cio.MultiWriter).Writer(name)
	return writer, found
}

func (cxt *ExecutionContext) AddLogger(name string, logger io.Writer) {
	cxt.loggers.(*cio.MultiWriter).AddWriter(name, logger)

	// Waiting until the first logger is attached before telling the Go log
	// system what the output is.
	if cxt.loggerRegistered == false {
		log.SetOutput(cxt.loggers)
		cxt.loggerRegistered = true
	}
}

func (cxt *ExecutionContext) RemoveLogger(name string) {
	cxt.loggers.(*cio.MultiWriter).RemoveWriter(name)
}

func (cxt *ExecutionContext) Log(prefix string, v ...interface{}) {
	tmpPrefix := log.Prefix()
	log.SetPrefix(prefix)
	log.Print(v...)
	log.SetPrefix(tmpPrefix)
}

func (cxt *ExecutionContext) Logf(prefix string, format string, v ...interface{}) {
	tmpPrefix := log.Prefix()
	log.SetPrefix(prefix)
	log.Printf(format, v...)
	log.SetPrefix(tmpPrefix)
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

	return newCxt
}
