package jason

type Null struct {
	Valid bool
}

func (j *Jason) Null() *Null {

	var valid bool

	// Check the type of this data
	switch j.data.(type) {
	case nil:
		valid = true
		break
	}

	n := new(Null)
	n.Valid = valid && j.exists // We also need to check that it actually exists here to separate nil and non-existing values

	return n
}

// Returns true if the object is actually an object
func (j *Jason) IsNull() bool {
	n := j.Null()
	return n.Valid
}
