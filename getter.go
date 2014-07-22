package cookoo

import (
	"reflect"
)

type Getter interface {
	Get(string, interface{}) interface{}
	Has(string) (interface{}, bool)
}

// GettableDS makes a KeyValueDatasource into a Getter.
//
// This is forward-compatibility code, and will be rendered unnecessary in
// Cookoo 2.x.
func GettableDS(ds KeyValueDatasource) Getter {
	return &gettableDatasource{ds}
}

// GettableCxt makes a Context into a Getter.
//
// This is forward-compatibility code, and will be rendered unnecessary in
// Cookoo 2.x.
func GettableCxt(cxt Context) Getter {
	return &gettableContext{cxt}
}

// GettableDatasource Makes a KeyValueDatasource match the Getter interface.
//
// In future versions of Cookoo, core Datasources will directly implement Getter.
type gettableDatasource struct {
	KeyValueDatasource
}

func (g *gettableDatasource) Get(key string, defaultVal interface{}) interface{} {
	ret := g.KeyValueDatasource.Value(key)
	if ret == nil || !reflect.ValueOf(ret).IsValid() {
		return defaultVal
	}
	return ret
}

func (g *gettableDatasource) Has(key string) (interface{}, bool) {
	ret := g.KeyValueDatasource.Value(key)
	if ret == nil || !reflect.ValueOf(ret).IsValid() {
		return nil, false
	}
	return ret, true
}

// GettableContext wraps a context and makes it a Getter.
// Since Context returns ContextValue objects, we have to write this stupid wrapper.
type gettableContext struct {
	Context
}

func (g *gettableContext) Get(key string, defaultVal interface{}) interface{} {
	return g.Context.Get(key, defaultVal)
}

func (g *gettableContext) Has(key string) (interface{}, bool) {
	return g.Context.Has(key)
}

// GetString is a convenience function for getting strings.
//
// This simplifies getting strings from a Context, a Params, or a
// GettableDatasource.
func GetString(key, defaultValue string, source Getter) string {
	out := source.Get(key, defaultValue)
	ret, ok := out.(string)
	if !ok {
		return defaultValue
	}
	return ret
}

func GetBool(key string, defaultValue bool, source Getter) bool {
	out := source.Get(key, defaultValue)
	ret, ok := out.(bool)
	if !ok {
		return defaultValue
	}
	return ret
}

func GetInt(key string, defaultValue int, source Getter) int {
	out := source.Get(key, defaultValue)
	ret, ok := out.(int)
	if !ok {
		return defaultValue
	}
	return ret
}

func GetInt64(key string, defaultValue int64, source Getter) int64 {
	out := source.Get(key, defaultValue)
	ret, ok := out.(int64)
	if !ok {
		return defaultValue
	}
	return ret
}

func GetInt32(key string, defaultValue int32, source Getter) int32 {
	out := source.Get(key, defaultValue)
	ret, ok := out.(int32)
	if !ok {
		return defaultValue
	}
	return ret
}

func GetUint64(key string, defaultVal uint64, source Getter) uint64 {
	out := source.Get(key, defaultVal)
	ret, ok := out.(uint64)
	if !ok {
		return defaultVal
	}
	return ret
}

func GetFloat64(key string, defaultVal float64, source Getter) float64 {
	out := source.Get(key, defaultVal)
	ret, ok := out.(float64)
	if !ok {
		return defaultVal
	}
	return ret
}

// HasString is a convenience function to perform Has() and return a string.
func HasString(key string, source Getter) (string, bool) {
	v, ok := source.Has(key)
	if !ok {
		return "", ok
	}
	strval, kk := v.(string)
	if !kk {
		return "", kk
	}
	return strval, kk
}

func HasBool(key string, defaultValue bool, source Getter) (bool, bool) {
	v, ok := source.Has(key)
	if !ok {
		return false, ok
	}
	strval, kk := v.(bool)
	if !kk {
		return false, kk
	}
	return strval, kk
}

func HasInt(key string, defaultValue int, source Getter) (int, bool) {
	v, ok := source.Has(key)
	if !ok {
		return 0, ok
	}
	val, kk := v.(int)
	if !kk {
		return 0, kk
	}
	return val, kk
}

func HasInt64(key string, defaultValue int64, source Getter) (int64, bool) {
	v, ok := source.Has(key)
	if !ok {
		return 0, ok
	}
	val, kk := v.(int64)
	if !kk {
		return 0, kk
	}
	return val, kk
}

func HasInt32(key string, defaultValue int32, source Getter) (int32, bool) {
	v, ok := source.Has(key)
	if !ok {
		return 0, ok
	}
	val, kk := v.(int32)
	if !kk {
		return 0, kk
	}
	return val, kk
}

func HasUint64(key string, defaultVal uint64, source Getter) (uint64, bool) {
	v, ok := source.Has(key)
	if !ok {
		return 0, ok
	}
	val, kk := v.(uint64)
	if !kk {
		return 0, kk
	}
	return val, kk
}

func HasFloat64(key string, defaultVal float64, source Getter) (float64, bool) {
	v, ok := source.Has(key)
	if !ok {
		return 0, ok
	}
	val, kk := v.(float64)
	if !kk {
		return 0, kk
	}
	return val, kk
}

// GetFromFirst gets the value from the first Getter that has the key.
//
// If no Getter has the key, the default value is returned, and the returned
// Getter is an instance of DefaultGetter.
func GetFromFirst(key string, defaultVal interface{}, sources ...Getter) (interface{}, Getter) {
	for _, s := range sources {
		val, ok := s.Has(key)
		if ok {
			return val, s
		}
	}

	return defaultVal, &DefaultGetter{defaultVal}
}

// DefaultGetter represents a Getter instance for a default value.
//
// A default getter always returns the given default value.
type DefaultGetter struct {
	val interface{}
}

func (e *DefaultGetter) Get(name string, value interface{}) interface{} {
	return e.val
}
func (e *DefaultGetter) Has(name string) (interface{}, bool) {
	return e.val, true
}
