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
type Array struct {
	Value
	slice []*Value // The formatted slice with typed values
	Valid bool
}

// Private bool
type Boolean struct {
	Value
	b     bool
	Valid bool
}

type Null struct {
	Value
	Valid bool
}

type Number struct {
	Value
	f     float64
	Valid bool
}

type Object struct {
	Value
	m     map[string]*Value // The formatted map with typed values
	Valid bool
}

func (v *Object) Map() map[string]*Value {
	return v.m
}

func (v *Array) Slice() []*Value {
	return v.slice
}

type String struct {
	Value
	Str   string
	Valid bool
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

// Marshal into bytes
func (j *Value) Marshal() ([]byte, error) {
	return json.Marshal(j.data)
}

// Private Get
func (j *Value) get(key string) (*Value, error) {

	// Assume this is an object
	obj := j.object()

	// Only continue if it really is an object
	if obj.Valid {
		child, ok := obj.Map()[key]
		if ok {
			return child, nil
		}
	}

	return nil, errors.New("could not get")

}

// Private to get path
func (j *Value) getPath(keys []string) (*Value, error) {
	current := j
	var err error
	for _, key := range keys {
		current, err = current.get(key)

		if err != nil {
			return nil, err
		}
	}
	return current, nil
}

// Get key or key path. Returns a new Value instance.
// Example: Get("address", "street")
func (j *Value) Get(keys ...string) (*Value, error) {
	return j.getPath(keys)
}

// Get Object at the path, and return error if it's not
// Can be useful in some cases

func (v *Value) GetObject(keys ...string) (*Object, error) {
	child, err := v.getPath(keys)

	if err != nil {
		return nil, err
	} else {

		obj, err := child.AsObject()

		if err != nil {
			return nil, err
		} else {
			return obj, nil
		}

	}

	return nil, nil
}

func (v *Value) GetString(keys ...string) (string, error) {
	child, err := v.getPath(keys)

	if err != nil {
		return "", err
	} else {

		obj, err := child.sstring()

		if err != nil {
			return "", err
		} else {
			return obj.AsString()
		}

	}

	return "", nil
}

func (v *Value) GetNumber(keys ...string) (float64, error) {
	child, err := v.getPath(keys)

	if err != nil {
		return 0, err
	} else {

		n, err := child.AsNumber()

		if err != nil {
			return 0, err
		} else {
			return n, nil
		}

	}

	return 0, nil
}

func (v *Value) GetBoolean(keys ...string) (bool, error) {
	child, err := v.getPath(keys)

	if err != nil {
		return false, err
	}

	return child.AsBoolean()
}

func (v *Value) GetArray(keys ...string) ([]*Value, error) {
	child, err := v.getPath(keys)

	if err != nil {
		return nil, err
	} else {

		return child.AsArray()

	}

	return nil, nil
}

/* // Not sure if we should keep this
// Determine if key path exists
func (j *Value) Has(keys ...string) bool {
	return j.getPath(keys).Exists()
}
*/

func (j *Value) null() (*Null, error) {

	var valid bool

	// Check the type of this data
	switch j.data.(type) {
	case nil:
		valid = true
		break
	}

	if valid {
		n := new(Null)
		n.Valid = valid && j.exists // We also need to check that it actually exists here to separate nil and non-existing values
		n.data = j.data
		return n, nil
	}

	return nil, errors.New("is not null")
}

// Returns true if the instance is actually a JSON null object.
func (j *Value) IsNull() bool {
	_, err := j.null()
	return err == nil
}

func (j *Value) array() *Array {

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
	var slice []*Value

	if valid {

		for _, element := range j.data.([]interface{}) {
			child := Value{element, true}
			slice = append(slice, &child)
		}
	}

	a.slice = slice
	a.data = j.data

	return a
}

// Returns the current data as an array of Jason values.
// Fallbacks on empty array
func (j *Value) AsArray() ([]*Value, error) {
	a := j.array()

	var err error

	if !a.Valid {
		err = errors.New("Is not an array")
	}

	return a.slice, err
}

func (j *Value) number() (*Number, error) {

	var valid bool

	// Check the type of this data
	switch j.data.(type) {
	case float64:
		valid = true
		break
	}

	if valid {
		n := new(Number)
		n.Valid = valid
		n.f = j.data.(float64)
		n.data = j.data
		return n, nil
	}

	return nil, errors.New("not a number")
}

func (j *Value) AsNumber() (float64, error) {
	n, err := j.number()

	if err != nil {
		return 0, err
	}

	return n.f, nil
}

// Returns true if the instance is actually a JSON bool.
func (j *Value) AsBoolean() (bool, error) {
	var valid bool

	// Check the type of this data
	switch j.data.(type) {
	case bool:
		valid = true
		break
	}

	if valid {
		return j.data.(bool), nil
	}

	return false, errors.New("no bool")
}

// Private object
func (j *Value) object() *Object {

	var valid bool

	// Check the type of this data
	switch j.data.(type) {
	case map[string]interface{}:
		valid = true
		break
	}

	obj := new(Object)
	obj.Valid = valid

	m := make(map[string]*Value)

	if valid {
		//obj.Map = j.data.(map[string]interface{})

		for key, element := range j.data.(map[string]interface{}) {
			m[key] = &Value{element, true}

		}
	}

	obj.data = j.data
	obj.m = m

	return obj
}

// Returns the current data as objects with string keys and Jason values.
// Fallbacks on empty map if invalid.
// Check IsObject() before using if you want to know.
func (j *Value) AsObject() (*Object, error) {
	obj := j.object()

	var err error

	if !obj.Valid {
		err = errors.New("Is not an object")
	}

	return obj, err
}

func (j *Value) sstring() (*String, error) {

	var valid bool

	// Check the type of this data
	switch j.data.(type) {
	case string:
		valid = true
		break
	}

	if valid {
		s := new(String)
		s.Valid = valid
		s.Str = j.data.(string)
		s.data = j.data
		return s, nil
	}

	return nil, errors.New("not a string")
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
	s, err := j.sstring()
	return s.String(), err

	/*
		s := j.sstring()

		var err error

		if !s.Valid {
			err = errors.New("Is not a string")
		}

		return s, err
	*/
}

// Returns true if the instance is actually a JSON string
func (v *Value) IsString() bool {
	_, err := v.sstring()
	return err == nil
}

// Used for logging
func (j *String) String() string {
	return j.data.(string)
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
