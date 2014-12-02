package jason

type Number struct {
	Float64 float64
	Valid   bool
}

func (j *Jason) Number() *Number {

	var valid bool

	// Check the type of this data
	switch j.data.(type) {
	case float64:
		valid = true
		break
	}

	n := new(Number)
	n.Valid = valid

	if valid {
		n.Float64 = j.data.(float64)
	}

	return n
}

// Returns true if the object is actually an object
func (j *Jason) IsNumber() bool {
	n := j.Number()
	return n.Valid
}
