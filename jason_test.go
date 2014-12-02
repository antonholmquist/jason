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
    }
  }`

	j, err := NewFromString(testJSON)

	assert.True(err == nil, "failed to create json from string")

	assert.True(j.Object().Valid, "the object should be valid")

	assert.True(j.Get("name").IsObject() == false, "name shoud not be an object")
	assert.True(j.Get("name").IsString() == true, "name shoud not be an object")
	assert.True(j.Get("name").String().String == "anton", "name shoud match")
	assert.True(j.Get("age").IsNumber() == true, "age should be a number")
	assert.True(j.Get("age").Number().Float64 == 29, "age mismatch")
	assert.True(j.Get("age").Exists(), "age should exist")
	assert.True(j.Get("age2").Exists() == false, "age2 should not exist")

	assert.True(j.Get("nothing").IsNull(), "nothing should be null")
	assert.True(j.Get("nothing2").IsNull() == false, "nothing2 fail")
	assert.True(j.Get("nothing").Exists(), "nothing should exist")
	assert.True(j.Get("nothing2").Exists() == false, "nothing2 should not exist")

	assert.True(j.Get("address").IsObject() == true, "address should be an object")
	assert.True(j.Get("address", "street").IsString() == true, "street should be a string")
	assert.True(j.Get("address", "street").String().String == "Street 42", "street mismatching")
	assert.True(j.Get("address", "street").Exists() == true, "street shoud exist")
	assert.True(j.Get("address", "street2").Exists() == false, "street should not exist")

	assert.True(j.Get("list").IsArray() == true, "list should be an array")

	for _, element := range j.Get("list").Array().Slice {
		// element is the element from someSlice for where we are
		log.Println("element: ", element)
	}
}
