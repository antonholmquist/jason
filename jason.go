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
func (v *Object) Map() map[string]*Value {
	return v.m
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

func objectFromValue(v *Value, err error) (*Object, error) {
	if err != nil {
		return nil, err
	} else {

		o, err := v.AsObject()

		if err != nil {
			return nil, err
		} else {
			return o, nil
		}

	}
}

func NewObjectFromString(s string) (*Object, error) {
	return objectFromValue(NewValueFromString(s))
}

func NewObjectFromBytes(b []byte) (*Object, error) {
	return objectFromValue(NewValueFromBytes(b))
}

func NewObjectFromReader(reader io.Reader) (*Object, error) {
	return objectFromValue(NewValueFromReader(reader))
}

// Marshal into bytes.
func (v *Value) Marshal() ([]byte, error) {
	return json.Marshal(v.data)
}

// Private Get
func (v *Value) get(key string) (*Value, error) {

	// Assume this is an object
	obj, err := v.AsObject()

	if err == nil {
		child, ok := obj.Map()[key]
		if ok {
			return child, nil
		} else {
			return nil, errors.New("key not found")
		}
	}

	return nil, err
}

// Private get path
func (v *Value) getPath(keys []string) (*Value, error) {
	current := v
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
func (v *Object) Get(keys ...string) (*Value, error) {
	return v.getPath(keys)
}

// Gets the value at key path and attempts to typecast the value into an object.
// Returns error if the value is not a json object.
func (v *Object) GetObject(keys ...string) (*Object, error) {
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
func (v *Object) GetString(keys ...string) (string, error) {
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
func (v *Object) GetNull(keys ...string) error {
	child, err := v.getPath(keys)

	if err != nil {
		return err
	}

	return child.AsNull()
}

// Gets the value at key path and attempts to typecast the value into a float64.
// Returns error if the value is not a json number.
func (v *Object) GetNumber(keys ...string) (float64, error) {
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
func (v *Object) GetBoolean(keys ...string) (bool, error) {
	child, err := v.getPath(keys)

	if err != nil {
		return false, err
	}

	return child.AsBoolean()
}

// Gets the value at key path and attempts to typecast the value into an array.
// Returns error if the value is not a json array.
func (v *Object) GetArray(keys ...string) ([]*Value, error) {
	child, err := v.getPath(keys)

	if err != nil {
		return nil, err
	} else {

		return child.AsArray()

	}

	return nil, nil
}

// Gets the value at key path and attempts to typecast the value into an array of objects.
// Returns error if the value is not a json array or if any of the contained objects are not objects.
func (v *Object) GetObjectArray(keys ...string) ([]*Object, error) {
	child, err := v.getPath(keys)

	if err != nil {
		return nil, err
	} else {

		array, err := child.AsArray()

		if err != nil {
			return nil, err
		} else {

			typedArray := make([]*Object, len(array))

			for index, arrayItem := range array {
				typedArrayItem, err := arrayItem.AsObject()

				if err != nil {
					return nil, err
				} else {
					typedArray[index] = typedArrayItem
				}

			}
			return typedArray, nil
		}
	}

	return nil, nil
}

// Gets the value at key path and attempts to typecast the value into an array of string.
// Returns error if the value is not a json array or if any of the contained objects are not strings.
func (v *Object) GetStringArray(keys ...string) ([]string, error) {
	child, err := v.getPath(keys)

	if err != nil {
		return nil, err
	} else {

		array, err := child.AsArray()

		if err != nil {
			return nil, err
		} else {

			typedArray := make([]string, len(array))

			for index, arrayItem := range array {
				typedArrayItem, err := arrayItem.AsString()

				if err != nil {
					return nil, err
				} else {
					typedArray[index] = typedArrayItem
				}

			}
			return typedArray, nil
		}
	}
	return nil, nil
}

// Gets the value at key path and attempts to typecast the value into an array of floats.
// Returns error if the value is not a json array or if any of the contained objects are not numbers.
func (v *Object) GetNumberArray(keys ...string) ([]float64, error) {
	child, err := v.getPath(keys)

	if err != nil {
		return nil, err
	} else {

		array, err := child.AsArray()

		if err != nil {
			return nil, err
		} else {

			typedArray := make([]float64, len(array))

			for index, arrayItem := range array {
				typedArrayItem, err := arrayItem.AsNumber()

				if err != nil {
					return nil, err
				} else {
					typedArray[index] = typedArrayItem
				}

			}
			return typedArray, nil
		}
	}
	return nil, nil
}

// Gets the value at key path and attempts to typecast the value into an array of bools.
// Returns error if the value is not a json array or if any of the contained objects are not booleans.
func (v *Object) GetBooleanArray(keys ...string) ([]bool, error) {
	child, err := v.getPath(keys)

	if err != nil {
		return nil, err
	} else {

		array, err := child.AsArray()

		if err != nil {
			return nil, err
		} else {

			typedArray := make([]bool, len(array))

			for index, arrayItem := range array {
				typedArrayItem, err := arrayItem.AsBoolean()

				if err != nil {
					return nil, err
				} else {
					typedArray[index] = typedArrayItem
				}

			}
			return typedArray, nil
		}
	}
	return nil, nil
}

// Gets the value at key path and attempts to typecast the value into an array of nulls.
// Returns length, or an error if the value is not a json array or if any of the contained objects are not nulls.
func (v *Object) GetNullArray(keys ...string) (int64, error) {
	child, err := v.getPath(keys)

	if err != nil {
		return 0, err
	} else {

		array, err := child.AsArray()

		if err != nil {
			return 0, err
		} else {

			var length int64 = 0

			for _, arrayItem := range array {
				err := arrayItem.AsNull()

				if err != nil {
					return 0, err
				} else {
					length++
				}

			}
			return length, nil
		}
	}

	return 0, nil
}

// Returns an error if the value is not actually null
func (v *Value) AsNull() error {
	var valid bool

	// Check the type of this data
	switch v.data.(type) {
	case nil:
		valid = v.exists // Valid only if j also exists, since other values could possibly also be nil
		break
	}

	if valid {
		return nil
	}

	return errors.New("is not null")

}

// Attempts to typecast the current value into an array.
// Returns error if the current value is not a json array.
func (v *Value) AsArray() ([]*Value, error) {
	var valid bool

	// Check the type of this data
	switch v.data.(type) {
	case []interface{}:
		valid = true
		break
	}

	// Unsure if this is a good way to use slices, it's probably not
	var slice []*Value

	if valid {

		for _, element := range v.data.([]interface{}) {
			child := Value{element, true}
			slice = append(slice, &child)
		}

		return slice, nil
	}

	return slice, errors.New("Not an array")

}

// Attempts to typecast the current value into a float64.
// Returns error if the current value is not a json number.
func (v *Value) AsNumber() (float64, error) {
	var valid bool

	// Check the type of this data
	switch v.data.(type) {
	case float64:
		valid = true
		break
	}

	if valid {
		return v.data.(float64), nil
	}

	return 0, errors.New("not a number")
}

// Attempts to typecast the current value into a bool.
// Returns error if the current value is not a json boolean.
func (v *Value) AsBoolean() (bool, error) {
	var valid bool

	// Check the type of this data
	switch v.data.(type) {
	case bool:
		valid = true
		break
	}

	if valid {
		return v.data.(bool), nil
	}

	return false, errors.New("no bool")
}

// Attempts to typecast the current value into an object.
// Returns error if the current value is not a json object.
func (v *Value) AsObject() (*Object, error) {

	var valid bool

	// Check the type of this data
	switch v.data.(type) {
	case map[string]interface{}:
		valid = true
		break
	}

	if valid {
		obj := new(Object)
		obj.valid = valid

		m := make(map[string]*Value)

		if valid {

			for key, element := range v.data.(map[string]interface{}) {
				m[key] = &Value{element, true}

			}
		}

		obj.data = v.data
		obj.m = m

		return obj, nil
	}

	return nil, errors.New("not an object")
}

// Attempts to typecast the current value into a string.
// Returns error if the current value is not a json string
func (v *Value) AsString() (string, error) {
	var valid bool

	// Check the type of this data
	switch v.data.(type) {
	case string:
		valid = true
		break
	}

	if valid {
		return v.data.(string), nil
	}

	return "", errors.New("not a string")
}

// Returns the value a json formatted string.
// Note: The method named String() is used by golang's log method for logging.
func (v *Value) String() string {
	f, err := json.Marshal(v.data)
	if err != nil {
		return err.Error()
	} else {
		return string(f)
	}
}
