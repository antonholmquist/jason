package jason

import (
	"encoding/json"
	"errors"
	"io"
)

// Value represents an arbitrary JSON value.
// It may contain a bool, number, string, object, array or null.
type Value struct {
	data   interface{}
	exists bool // Used to separate nil and non-existing values
}

// Object represents an object JSON value.
// The underlying golang map can be accessed with Map.
// It is needed when iterating through the values of the object.
type Object struct {
	Value
	Map   map[string]*Value // The formatted map with typed values
	Valid bool
}

// Create a new Value from a io.reader.
// Useful for parsing the body of a net/http response.
// Returns an error if something went wrong.
// Example: NewFromReader(res.Body)
func NewValueFromReader(reader io.Reader) (*Value, error) {
	j := new(Value)
	d := json.NewDecoder(reader)
	err := d.Decode(&j.data)
	return j, err
}

// Create a new Value from bytes
// Returns an error if the bytes couldn't be parsed.
func NewValueFromBytes(b []byte) (*Value, error) {
	j := new(Value)
	err := json.Unmarshal(b, &j.data)
	return j, err
}

// Create a new Value from a string
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
		child, ok := obj.Map[key]
		if ok {
			return child, nil
		}
	}

	return nil, errors.New("could not get")

}

// Private get path
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

// Get Value at key path. Returns a Value or error.
// Example: Get("address", "street")
func (j *Value) Get(keys ...string) (*Value, error) {
	return j.getPath(keys)
}

// Get Object at key path. Returns an Object, or error.
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

// Get string at key path. Returns a string primitive, or error.
func (v *Value) GetString(keys ...string) (string, error) {
	child, err := v.getPath(keys)

	if err != nil {
		return "", err
	} else {

		return child.AsString()

	}

	return "", nil
}

// Get number at key path. Returns a float64, or error.
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

// Get number at key path. Returns a bool, or error.
func (v *Value) GetBoolean(keys ...string) (bool, error) {
	child, err := v.getPath(keys)

	if err != nil {
		return false, err
	}

	return child.AsBoolean()
}

// Get array at key path. Returns an array of Values, or error.
func (v *Value) GetArray(keys ...string) ([]*Value, error) {
	child, err := v.getPath(keys)

	if err != nil {
		return nil, err
	} else {

		return child.AsArray()

	}

	return nil, nil
}

// Returns true if the instance is actually a JSON null object.
func (j *Value) IsNull() bool {
	var valid bool

	// Check the type of this data
	switch j.data.(type) {
	case nil:
		valid = j.exists // Valid only if j also exists, since other values could possibly also be nil
		break
	}

	return valid

}

// Returns the current Value as an array of Values.
// Returns an error if the current Value is not an Array.
func (j *Value) AsArray() ([]*Value, error) {
	var valid bool

	// Check the type of this data
	switch j.data.(type) {
	case []interface{}:
		valid = true
		break
	}

	// Unsure if this is a good way to use slices, it's probably not
	var slice []*Value

	if valid {

		for _, element := range j.data.([]interface{}) {
			child := Value{element, true}
			slice = append(slice, &child)
		}

		return slice, nil
	}

	return slice, errors.New("Not an array")

}

// Returns the current Value as a float64.
// Returns an error if the current Value is not a json number.
func (j *Value) AsNumber() (float64, error) {
	var valid bool

	// Check the type of this data
	switch j.data.(type) {
	case float64:
		valid = true
		break
	}

	if valid {
		return j.data.(float64), nil
	}

	return 0, errors.New("not a number")
}

// Returns the current Value as a bool.
// Returns an error if the current Value is not a json boolean.
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
	obj.Map = m

	return obj
}

// Returns the current data as objects with string keys and Value values.
// Returns an error if the current Value is not a json object
func (j *Value) AsObject() (*Object, error) {
	obj := j.object()

	var err error

	if !obj.Valid {
		err = errors.New("Is not an object")
	}

	return obj, err
}

// Returns the current data as a golang string
// Returns an error if the current Value is not a json string
func (j *Value) AsString() (string, error) {
	var valid bool

	// Check the type of this data
	switch j.data.(type) {
	case string:
		valid = true
		break
	}

	if valid {
		return j.data.(string), nil
	}

	return "", errors.New("not a string")
}

// The method named String() is used by golang's log method for logging
// It returns the Value as object
func (j *Value) String() string {
	f, err := json.Marshal(j.data)
	if err != nil {
		return err.Error()
	} else {
		return string(f)
	}
}
