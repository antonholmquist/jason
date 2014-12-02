package jason

import (
	"encoding/json"
	"io"
	_ "log"
)

type Jason struct {
	data   interface{}
	exists bool // Used to separate nil and non-existing values
}

func NewFromReader(reader io.Reader) (*Jason, error) {
	j := new(Jason)
	d := json.NewDecoder(reader)
	err := d.Decode(&j.data)
	return j, err
}

func NewFromString(s string) (*Jason, error) {
	j := new(Jason)
	b := []byte(s)
	err := json.Unmarshal(b, &j.data)
	return j, err
}

func (j *Jason) Marshal([]byte, error) {
	json.Marshal(j.data)
}

// Returns true if this key exists
func (j *Jason) Exists() bool {
	return j.exists
}

// Get
func (j *Jason) Get(key string) *Jason {

	// Assume this is an object
	obj := j.Object()

	// Only continue if it really is an object
	if obj.Valid {
		childData, ok := obj.Map[key]
		if ok {
			return &Jason{childData, true}
		}
	}

	return &Jason{nil, false}

}
