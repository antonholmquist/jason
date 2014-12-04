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

	a, err := j.GetObject("address")
	assert.True(a != nil && err == nil, "failed to create json from string")

	assert.True(err == nil, "failed to create json from string")

	s, err := j.GetString("name")

	assert.True(s.String() == "anton" && err == nil, "name should be a string")
	//assert.True(j.Get("name").IsObject() == false, "name should not be an object")

	assert.True(j.object().Valid, "the object should be valid")

	//assert.True(j.Has("name") == true, "has name")
	//assert.True(j.Has("name2") == false, "do not have name2")

	s, err = j.GetString("name")
	assert.True(s.String() == "anton" && err == nil, "name shoud match")

	s, err = j.GetString("address", "street")
	assert.True(s.String() == "Street 42" && err == nil, "street shoud match")
	//log.Println("s: ", s.String())

	_, err = j.GetNumber("age")
	assert.True(err == nil, "age should be a number")

	n, err := j.GetNumber("age")
	assert.True(n.Float64() == 29.0 && err == nil, "age mismatch")

	age, err := j.Get("age")
	assert.True(age != nil && err == nil, "age should exist")

	age2, err := j.Get("age2")
	assert.True(age2 == nil && err != nil, "age2 should not exist")

	address, err := j.GetObject("address")
	assert.True(address != nil && err == nil, "address should be an object")

	//log.Println("address: ", address)

	s, err = address.GetString("street")

	addressAsString, err := j.GetString("address")
	assert.True(addressAsString == nil && err != nil, "address should not be an string")

	s, err = j.GetString("address", "street")
	assert.True(s.String() == "Street 42" && err == nil, "street mismatching")

	s, err = j.GetString("address", "name2")
	assert.True(s == nil && err != nil, "nonexistent string fail")

	b, err := j.GetBoolean("true")
	assert.True(b.Boolean() == true && err == nil, "bool true test")

	b, err = j.GetBoolean("false")
	assert.True(b.Boolean() == false && err == nil, "bool false test")

	b, err = j.GetBoolean("invalid_field")
	assert.True(b == nil && err != nil, "bool invalid test")

	list, err := j.GetArray("list")
	assert.True(list != nil && err == nil, "list should be an array")

	list2, err := j.GetArray("list2")
	assert.True(list2 != nil && err == nil, "list2 should be an array")

	list2Array, err := j.GetArray("list2")
	assert.True(err == nil, "List2 should not return error on AsArray")
	assert.True(len(list2Array.Slice()) == 2, "List2 should should have length 2")

	for _, element := range list2Array.Slice() {
		//assert.True(element.IsObject() == true, "first fail")

		s, err = element.GetString("street")
		assert.True(s.String() == "Street 42" && err == nil, "second fail")
	}

	obj, err := j.GetObject("country")
	assert.True(obj != nil && err == nil, "country should not return error on AsObject")
	for key, value := range obj.Map() {

		assert.True(key == "name", "country name key incorrect")

		s, err = value.AsString()
		assert.True(s.String() == "Sweden" && err == nil, "country name should be Sweden")
	}
}
