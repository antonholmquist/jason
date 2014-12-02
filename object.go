package jason

type Object struct {
	Map   map[string]interface{}
	Valid bool
}

// Get the current instance as an object
// Will always return a value object.
// To check if it's a valid json object, check Valid flag of the returned object.
func (j *Jason) Object() *Object {

	var valid bool

	// Check the type of this data
	switch j.data.(type) {
	case map[string]interface{}:
		valid = true
		break
	}

	obj := new(Object)
	obj.Valid = valid

	if valid {
		obj.Map = j.data.(map[string]interface{})
	}

	return obj
}

// Returns true if the object is actually an object
func (j *Jason) IsObject() bool {
	obj := j.Object()
	return obj.Valid
}
