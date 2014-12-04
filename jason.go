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

type Object struct {
	Value
	Map   map[string]*Value // The formatted map with typed values
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
		child, ok := obj.Map[key]
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

		return child.AsString()

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

// Returns the current data as an array of Jason values.
// Fallbacks on empty array
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
	obj.Map = m

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

// Returns the current data as string. Fallbacks on empty string if invalid.
// Check IsString() before using if you want to know.
// It's good to use this same since String() conflicts with log default method
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
// We therefore send the json
func (j *Value) String() string {
	f, err := json.Marshal(j.data)
	if err != nil {
		return err.Error()
	} else {
		return string(f)
	}
}
