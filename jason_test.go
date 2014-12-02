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

	assert.True(true, "this should be true")
	assert.True(false, "this should be true")

	log.Println("FirstTest")
}
