package object

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironment(t *testing.T) {
	e := NewEnvironment()
	obj := &String{Value: "value"}
	e.Set("test", obj)

	result, ok := e.Get("test")

	assert.True(t, ok)
	assert.Equal(t, "value", result.Inspect())
}
