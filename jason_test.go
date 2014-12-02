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
    "address": {
      "street": "Street 42",
      "city": "Stockholm"
    }
  }`

	log.Println("FirstTest: ", testJSON)

	j, err := NewFromString(testJSON)

	log.Println("data: ", j.data)

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

	assert.True(j.Get("list").IsArray() == true, "list should be an array")

}
