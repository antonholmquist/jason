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

	j, err := NewFromString(testJSON)

	assert.True(err == nil, "failed to create json from string")

	assert.True(j.Get("name").IsString() == true, "name should be a string")
	assert.True(j.Get("name").IsObject() == false, "name should not be an object")

	assert.True(j.object().Valid, "the object should be valid")

	assert.True(j.Has("name") == true, "has name")
	assert.True(j.Has("name2") == false, "do not have name2")
	assert.True(j.Get("name").String() == "anton", "name shoud match")
	assert.True(j.Get("age").IsNumber() == true, "age should be a number")
	assert.True(j.Get("age").Number() == 29.0, "age mismatch")
	assert.True(j.Get("age").Exists(), "age should exist")
	assert.True(j.Get("age2").Exists() == false, "age2 should not exist")

	assert.True(j.Get("nothing").IsNull(), "nothing should be null")
	assert.True(j.Get("nothing2").IsNull() == false, "nothing2 fail")
	assert.True(j.Get("nothing").Exists(), "nothing should exist")
	assert.True(j.Get("nothing2").Exists() == false, "nothing2 should not exist")

	assert.True(j.Get("address").IsObject() == true, "address should be an object")
	assert.True(j.Get("address", "street").IsString() == true, "street should be a string")
	assert.True(j.Get("address", "street").String() == "Street 42", "street mismatching")
	assert.True(j.Get("address", "street").Exists() == true, "street shoud exist")
	assert.True(j.Get("address", "street2").Exists() == false, "street should not exist")

	assert.True(j.Get("true").IsBoolean(), "true test")
	assert.True(j.Get("false").IsBoolean(), "true test")
	assert.True(j.Get("true").Boolean() == true, "true test")
	assert.True(j.Get("false").Boolean() == false, "true test")

	assert.True(j.Get("list").IsArray() == true, "list should be an array")
	assert.True(j.Get("list2").IsArray() == true, "list2 should be an array")
	assert.True(j.Get("list2") != nil, "list2 should exist")

	for _, element := range j.Get("list2").Array() {
		assert.True(element.IsObject() == true, "first fail")
		assert.True(element.Get("street").String() == "Street 42", "second fail")
	}

	for key, value := range j.Get("country").Object() {

		assert.True(key == "name", "country name key incorrect")
		assert.True(value.IsString(), "country name should be a string")
		assert.True(value.String() == "Sweden", "country name should be Sweden")
	}
}
