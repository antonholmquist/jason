package jason

import (
	"encoding/json"
	"io"
)

type Jason struct {
	data   interface{}
	exists bool // Used to separate nil and non-existing values
}

// Create a new instance from a io.reader
func NewFromReader(reader io.Reader) (*Jason, error) {
	j := new(Jason)
	d := json.NewDecoder(reader)
	err := d.Decode(&j.data)
	return j, err
}

// Create a new instance from bytes
func NewFromBytes(b []byte) (*Jason, error) {
	j := new(Jason)
	err := json.Unmarshal(b, &j.data)
	return j, err
}

// Create a new instance from a string
func NewFromString(s string) (*Jason, error) {
	b := []byte(s)
	return NewFromBytes(b)
}

func (j *Jason) Marshal([]byte, error) {
	json.Marshal(j.data)
}

// Returns true if this key exists
func (j *Jason) Exists() bool {
	return j.exists
}

// Private Get
func (j *Jason) get(key string) *Jason {

	// Assume this is an object
	obj := j.object()

	// Only continue if it really is an object
	if obj.Valid {
		child, ok := obj.Map[key]
		if ok {
			return child
		}
	}

	return &Jason{nil, false}

}

// Get key
// Example: Get("address", "street")
func (j *Jason) Get(args ...string) *Jason {
	current := j
	for _, key := range args {
		current = current.get(key)
	}
	return current
}
