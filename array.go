package jason

type Array struct {
	Slice []interface{}
	Valid bool
}

func (j *Jason) Array() *Array {

	var valid bool

	// Check the type of this data
	switch j.data.(type) {
	case []interface{}:
		valid = true
		break
	}

	a := new(Array)
	a.Valid = valid

	if valid {
		a.Slice = j.data.([]interface{})
	}

	return a
}

// Returns true if the array is actually an array
func (j *Jason) IsArray() bool {
	a := j.Array()
	return a.Valid
}
