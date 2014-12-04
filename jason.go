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

// Object represents an object JSON object.
// It inherets from Value but with an additional method to access
// a map representation of it's content. It's useful when iterating.
type Object struct {
	Value
	m     map[string]*Value
	valid bool
}

// Returns the golang map.
// Needed when iterating through the values of the object.
func (j *Object) Map() map[string]*Value {
	return j.m
}

// Creates a new value from an io.reader.
// Returns an error if the reader does not contain valid json.
// Useful for parsing the body of a net/http response.
// Example: NewFromReader(res.Body)
func NewValueFromReader(reader io.Reader) (*Value, error) {
	j := new(Value)
	d := json.NewDecoder(reader)
	err := d.Decode(&j.data)
	return j, err
}

// Creates a new value from bytes.
// Returns an error if the bytes are not valid json.
func NewValueFromBytes(b []byte) (*Value, error) {
	j := new(Value)
	err := json.Unmarshal(b, &j.data)
	return j, err
}

// Creates a new value from a string.
// Returns an error if the string is not valid json.
func NewValueFromString(s string) (*Value, error) {
	b := []byte(s)
	return NewValueFromBytes(b)
}

// Marshal into bytes.
func (j *Value) Marshal() ([]byte, error) {
	return json.Marshal(j.data)
}

// Private Get
func (j *Value) get(key string) (*Value, error) {

	// Assume this is an object
	obj := j.object()

	// Only continue if it really is an object
	if obj.valid {
		child, ok := obj.Map()[key]
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

// Gets the value at key path.
// Returns error if the value does not exist.
// Example: Get("address", "street")
func (j *Value) Get(keys ...string) (*Value, error) {
	return j.getPath(keys)
}

// Gets the value at key path and attempts to typecast the value into an object.
// Returns error if the value is not a json object.
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

// Gets the value at key path and attempts to typecast the value into a string.
// Returns error if the value is not a json string.
func (v *Value) GetString(keys ...string) (string, error) {
	child, err := v.getPath(keys)

	if err != nil {
		return "", err
	} else {
		return child.AsString()
	}

	return "", nil
}

// Gets the value at key path and attempts to typecast the value into null.
// Returns error if the value is not json null.
func (v *Value) GetNull(keys ...string) error {
	child, err := v.getPath(keys)

	if err != nil {
		return err
	}

	return child.AsNull()
}

// Gets the value at key path and attempts to typecast the value into a float64.
// Returns error if the value is not a json number.
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

// Gets the value at key path and attempts to typecast the value into a bool.
// Returns error if the value is not a json boolean.
func (v *Value) GetBoolean(keys ...string) (bool, error) {
	child, err := v.getPath(keys)

	if err != nil {
		return false, err
	}

	return child.AsBoolean()
}

// Gets the value at key path and attempts to typecast the value into an array.
// Returns error if the value is not a json array.
func (v *Value) GetArray(keys ...string) ([]*Value, error) {
	child, err := v.getPath(keys)

	if err != nil {
		return nil, err
	} else {

		return child.AsArray()

	}

	return nil, nil
}

// Returns an error if the value is not actually null
func (j *Value) AsNull() error {
	var valid bool

	// Check the type of this data
	switch j.data.(type) {
	case nil:
		valid = j.exists // Valid only if j also exists, since other values could possibly also be nil
		break
	}

	if valid {
		return nil
	}

	return errors.New("is not null")

}

// Attempts to typecast the current value into an array.
// Returns error if the current value is not a json array.
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

// Attempts to typecast the current value into a float64.
// Returns error if the current value is not a json number.
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

// Attempts to typecast the current value into a bool.
// Returns error if the current value is not a json boolean.
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
	obj.valid = valid

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

// Attempts to typecast the current value into an object.
// Returns error if the current value is not a json object.
func (j *Value) AsObject() (*Object, error) {
	obj := j.object()

	var err error

	if !obj.valid {
		err = errors.New("Is not an object")
	}

	return obj, err
}

// Attempts to typecast the current value into a string.
// Returns error if the current value is not a json string
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

// Returns the value a json formatted string.
// Note: The method named String() is used by golang's log method for logging.
func (j *Value) String() string {
	f, err := json.Marshal(j.data)
	if err != nil {
		return err.Error()
	} else {
		return string(f)
	}
}
