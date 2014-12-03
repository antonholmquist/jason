package jason

import (
	"log"
	"testing"
)

type Assert struct {
	T *testing.T
}

func NewAssert(t *testing.T) *Assert {
	return &Assert{
		T: t,
	}
}

func (assert *Assert) True(value bool, message string) {
	if value == false {
		log.Panicln("Assert: ", message)
	}
}

func TestFirst(t *testing.T) {

	assert := NewAssert(t)

	testJSON := `{
    "name": "anton",
    "age": 29,
    "nothing": null,
    "true": true,
    "false": false,
    "list": [
      "first",
      "second"
    ],
    "list2": [
      {
        "street": "Street 42",
        "city": "Stockholm"
      },
      {
        "street": "Street 42",
        "city": "Stockholm"
      }
    ],
    "address": {
      "street": "Street 42",
      "city": "Stockholm"
    },
    "country": {
      "name": "Sweden"
    }
  }`

	j, err := NewValueFromString(testJSON)

	assert.True(err == nil, "failed to create json from string")

	s, err := j.Get("name").AsString()
	assert.True(s == "anton" && err == nil, "name should be a string")
	//assert.True(j.Get("name").IsObject() == false, "name should not be an object")

	assert.True(j.object().Valid, "the object should be valid")

	//assert.True(j.Has("name") == true, "has name")
	//assert.True(j.Has("name2") == false, "do not have name2")

	s, err = j.Get("name").AsString()
	assert.True(s == "anton" && err == nil, "name shoud match")

	_, err = j.Get("age").AsNumber()
	assert.True(err == nil, "age should be a number")

	n, err := j.Get("age").AsNumber()
	assert.True(n == 29.0 && err == nil, "age mismatch")
	assert.True(j.Get("age").Exists(), "age should exist")
	assert.True(j.Get("age2").Exists() == false, "age2 should not exist")

	assert.True(j.Get("nothing").IsNull(), "nothing should be null")
	assert.True(j.Get("nothing2").IsNull() == false, "nothing2 fail")
	assert.True(j.Get("nothing").Exists(), "nothing should exist")
	assert.True(j.Get("nothing2").Exists() == false, "nothing2 should not exist")

	address, err := j.Get("address").AsObject()
	assert.True(address != nil && err == nil, "address should be an object")

	addressAsString, err := j.Get("address").AsString()
	assert.True(addressAsString == "" && err != nil, "address should not be an string")

	s, err = j.Get("address", "street").AsString()
	assert.True(s == "Street 42" && err == nil, "street mismatching")

	s, err = j.Get("address", "name2").AsString()
	assert.True(s == "" && err != nil, "nonexistent string fail")

	assert.True(j.Get("address", "street").Exists() == true, "street shoud exist")
	assert.True(j.Get("address", "street2").Exists() == false, "street should not exist")

	b, err := j.Get("true").AsBoolean()
	assert.True(b == true && err == nil, "bool true test")

	b, err = j.Get("false").AsBoolean()
	assert.True(b == false && err == nil, "bool false test")

	b, err = j.Get("invalid_field").AsBoolean()
	assert.True(b == false && err != nil, "bool invalid test")

	list, err := j.Get("list").AsArray()
	assert.True(list != nil && err == nil, "list should be an array")

	list2, err := j.Get("list2").AsArray()
	assert.True(list2 != nil && err == nil, "list2 should be an array")

	list2Array, err := j.Get("list2").AsArray()
	assert.True(err == nil, "List2 should not return error on AsArray")
	assert.True(len(list2Array) == 2, "List2 should should have length 2")

	for _, element := range list2Array {
		//assert.True(element.IsObject() == true, "first fail")

		s, err = element.Get("street").AsString()
		assert.True(s == "Street 42" && err == nil, "second fail")
	}

	obj, err := j.Get("country").AsObject()
	assert.True(obj != nil && err == nil, "country should not return error on AsObject")
	for key, value := range obj {

		assert.True(key == "name", "country name key incorrect")

		s, err = value.AsString()
		assert.True(s == "Sweden" && err == nil, "country name should be Sweden")
	}
}
