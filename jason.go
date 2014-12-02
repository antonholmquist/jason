package jason

import (
	"encoding/json"
	"io"
)

type Jason struct {
	data interface{}
}

func NewFromReader(reader io.Reader) (*Jason, error) {
	j := new(Jason)
	d := json.NewDecoder(reader)
	d.UseNumber()
	err := d.Decode(&j.data)
	return j, err
}

func (j *Jason) Marshal([]byte, error) {
	json.Marshal(j.data)
}
