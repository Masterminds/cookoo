package cookoo

type Params map[string]interface{}

// Check if a parameter exists, and return it if found.
func (p Params) Has(name string) (value interface{}, ok bool) {
	value, ok = p[name]
	return
}

// Get a parameter value, or return the default value.
func (p Params) Get(name string, defaultValue interface{}) interface{} {
	val, ok := p.Has(name)
	if ok {
		return val
	}
	return defaultValue
}

// Require that a given list of parameters are present.
// If they are all present, ok = true. Otherwise, ok = false and the 
// `missing` array contains a list of missing params.
func (p Params) Requires(paramNames ...string) (ok bool, missing []string) {
	missing = make([]string, 0, len(p))
	for _, val := range paramNames {
		_, ok := p[val]
		if !ok {
			missing = append(missing, val)
		}
	}
	ok = len(missing) == 0
	return
}

// Given a name and a validation function, return a valid value.
// If the value is not valid, ok = false.
func (p Params) Validate(name string, validator func(interface{})bool) (value interface{}, ok bool) {
	value, ok = p[name]
	if !ok {
		return
	}

	if !validator(value.(interface{})) {
		// XXX: For safety, we set a failed value to nil.
		value = nil
		ok = false
	}
	return
}
