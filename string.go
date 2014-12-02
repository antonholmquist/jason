package jason

type String struct {
	String string
	Valid  bool
}

func (j *Jason) String() *String {

	var valid bool

	// Check the type of this data
	switch j.data.(type) {
	case string:
		valid = true
		break
	}

	s := new(String)
	s.Valid = valid

	if valid {
		s.String = j.data.(string)
	}

	return s
}

// Returns true if the object is actually an object
func (j *Jason) IsString() bool {
	s := j.String()
	return s.Valid
}
