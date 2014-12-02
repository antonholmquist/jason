package jason

import (
	_ "log"
)

type Object struct {
	Map   map[string]*Jason
	Valid bool
}

// Get the current instance as an object
// Will always return a value object.
// To check if it's a valid json object, check Valid flag of the returned object.
func (j *Jason) object() *Object {

	var valid bool

	// Check the type of this data
	switch j.data.(type) {
	case map[string]interface{}:
		valid = true
		break
	}

	obj := new(Object)
	obj.Valid = valid

	m := make(map[string]*Jason)

	if valid {
		//obj.Map = j.data.(map[string]interface{})

		for key, element := range j.data.(map[string]interface{}) {
			m[key] = &Jason{element, true}
		}
	}

	obj.Map = m

	return obj
}

func (j *Jason) Object() map[string]*Jason {
	obj := j.object()
	return obj.Map
}

// Returns true if the object is actually an object
func (j *Jason) IsObject() bool {
	obj := j.object()
	return obj.Valid
}
