package jason

type Array struct {
	Slice []*Jason
	Valid bool
}

func (j *Jason) array() *Array {

	var valid bool

	// Check the type of this data
	switch j.data.(type) {
	case []interface{}:
		valid = true
		break
	}

	a := new(Array)
	a.Valid = valid

	// Unsure if this is a good way to use slices, it's probably not
	var slice []*Jason

	if valid {

		for _, element := range j.data.([]interface{}) {
			child := Jason{element, true}
			slice = append(slice, &child)
		}
	}

	a.Slice = slice

	return a
}

func (j *Jason) Array() []*Jason {
	a := j.array()
	return a.Slice
}

// Returns true if the array is actually an array
func (j *Jason) IsArray() bool {
	a := j.array()
	return a.Valid
}
