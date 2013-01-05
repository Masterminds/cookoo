package cookoo

type Params map[string]interface{}

func (p Params) Has(name string) (value interface{}, ok bool) {
	value, ok = p[name]
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
