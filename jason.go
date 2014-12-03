package jason

import (
	"encoding/json"
	"errors"
	"io"
)

type Value struct {
	data   interface{}
	exists bool // Used to separate nil and non-existing values
}

// Private array
type jArray struct {
	Slice []*Value
	Valid bool
}

// Private bool
type jBool struct {
	Bool  bool
	Valid bool
}

type jNull struct {
	Valid bool
}

type jNumber struct {
	Float64 float64
	Valid   bool
}

type jObject struct {
	Map   map[string]*Value
	Valid bool
}

type jString struct {
	String string
	Valid  bool
}

// Create a new instance from a io.reader.
// Useful for parsing the body of a net/http response.
// Example: NewFromReader(res.Body)
func NewValueFromReader(reader io.Reader) (*Value, error) {
	j := new(Value)
	d := json.NewDecoder(reader)
	err := d.Decode(&j.data)
	return j, err
}

// Create a new instance from bytes
// Returns an error if the bytes couldn't be parsed.
func NewValueFromBytes(b []byte) (*Value, error) {
	j := new(Value)
	err := json.Unmarshal(b, &j.data)
	return j, err
}

// Create a new instance from a string
// Returns an error if the string couldn't be parsed.
func NewValueFromString(s string) (*Value, error) {
	b := []byte(s)
	return NewValueFromBytes(b)
}

// Returns true if this key exists
// Example: j.Get("address").Exists()
func (j *Value) Exists() bool {
	return j.exists
}

// Marshal into bytes
func (j *Value) Marshal() ([]byte, error) {
	return json.Marshal(j.data)
}

// Private Get
func (j *Value) get(key string) *Value {

	// Assume this is an object
	obj := j.object()

	// Only continue if it really is an object
	if obj.Valid {
		child, ok := obj.Map[key]
		if ok {
			return child
		}
	}

	return &Value{nil, false}

}

// Private to get path
func (j *Value) getPath(keys []string) *Value {
	current := j
	for _, key := range keys {
		current = current.get(key)
	}
	return current
}

// Get key or key path. Returns a new Value instance.
// Example: Get("address", "street")
func (j *Value) Get(keys ...string) *Value {
	return j.getPath(keys)
}

/* // Not sure if we should keep this
// Determine if key path exists
func (j *Value) Has(keys ...string) bool {
	return j.getPath(keys).Exists()
}
*/

func (j *Value) null() *jNull {

	var valid bool

	// Check the type of this data
	switch j.data.(type) {
	case nil:
		valid = true
		break
	}

	n := new(jNull)
	n.Valid = valid && j.exists // We also need to check that it actually exists here to separate nil and non-existing values

	return n
}

// Returns true if the instance is actually a JSON null object.
func (j *Value) IsNull() bool {
	n := j.null()
	return n.Valid
}

func (j *Value) array() *jArray {

	var valid bool

	// Check the type of this data
	switch j.data.(type) {
	case []interface{}:
		valid = true
		break
	}

	a := new(jArray)
	a.Valid = valid

	// Unsure if this is a good way to use slices, it's probably not
	var slice []*Value

	if valid {

		for _, element := range j.data.([]interface{}) {
			child := Value{element, true}
			slice = append(slice, &child)
		}
	}

	a.Slice = slice

	return a
}

// Returns the current data as an array of Jason values.
// Fallbacks on empty array
// Check IsArray() before using if you want to know.
func (j *Value) AsArray() ([]*Value, error) {
	a := j.array()

	var err error

	if !a.Valid {
		err = errors.New("Is not an array")
	}

	return a.Slice, err
}

func (j *Value) IsArray() bool {
	a := j.array()
	return a.Valid
}

func (j *Value) number() *jNumber {

	var valid bool

	// Check the type of this data
	switch j.data.(type) {
	case float64:
		valid = true
		break
	}

	n := new(jNumber)
	n.Valid = valid

	if valid {
		n.Float64 = j.data.(float64)
	}

	return n
}

func (j *Value) AsNumber() (float64, error) {
	n := j.number()

	var err error

	if !n.Valid {
		err = errors.New("Is not a number")
	}

	return n.Float64, err
}

func (j *Value) IsNumber() bool {
	n := j.number()
	return n.Valid
}

// Returns the same as Number()
func (j *Value) AsFloat64() (float64, error) {
	return j.AsNumber()
}

// Returns the Number() converted to an int64
func (j *Value) AsInt64() (int64, error) {
	f, err := j.AsNumber()
	return int64(f), err
}

// Private
func (j *Value) boolean() *jBool {

	var valid bool

	// Check the type of this data
	switch j.data.(type) {
	case bool:
		valid = true
		break
	}

	b := new(jBool)
	b.Valid = valid

	if valid {
		b.Bool = j.data.(bool)
	}

	return b
}

// Returns true if the instance is actually a JSON bool.
func (j *Value) AsBoolean() (bool, error) {
	b := j.boolean()
	var err error

	if !b.Valid {
		err = errors.New("Is not a bool")
	}

	return b.Bool, err
}

func (v *Value) IsBoolean() bool {
	b := v.boolean()
	return b.Valid
}

// Private object
func (j *Value) object() *jObject {

	var valid bool

	// Check the type of this data
	switch j.data.(type) {
	case map[string]interface{}:
		valid = true
		break
	}

	obj := new(jObject)
	obj.Valid = valid

	m := make(map[string]*Value)

	if valid {
		//obj.Map = j.data.(map[string]interface{})

		for key, element := range j.data.(map[string]interface{}) {
			m[key] = &Value{element, true}
		}
	}

	obj.Map = m

	return obj
}

// Returns the current data as objects with string keys and Jason values.
// Fallbacks on empty map if invalid.
// Check IsObject() before using if you want to know.
func (j *Value) AsObject() (map[string]*Value, error) {
	obj := j.object()

	var err error

	if !obj.Valid {
		err = errors.New("Is not an object")
	}

	return obj.Map, err
}

func (j *Value) sstring() *jString {

	var valid bool

	// Check the type of this data
	switch j.data.(type) {
	case string:
		valid = true
		break
	}

	s := new(jString)
	s.Valid = valid

	if valid {
		s.String = j.data.(string)
	}

	return s
}

// Returns true if the instance is actually a JSON object
func (v *Value) IsObject() bool {
	obj := v.object()
	return obj.Valid
}

// Returns the current data as string. Fallbacks on empty string if invalid.
// Check IsString() before using if you want to know.
// It's good to use this same since String() conflicts with log default method
func (j *Value) AsString() (string, error) {

	s := j.sstring()

	var err error

	if !s.Valid {
		err = errors.New("Is not a string")
	}

	return s.String, err
}

// Returns true if the instance is actually a JSON string
func (v *Value) IsString() bool {
	s := v.sstring()
	return s.Valid
}

// Used for logging
func (j *Value) String() string {
	f, err := json.Marshal(j.data)
	if err != nil {
		return err.Error()
	} else {
		return string(f)
	}
}
