package jason

import ()

type String struct {
	String string
	Valid  bool
}

func (j *Jason) sstring() *String {

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

func (j *Jason) String() string {
	s := j.sstring()
	return s.String
}

// Returns true if the object is actually an object
func (j *Jason) IsString() bool {
	s := j.sstring()
	return s.Valid
}
