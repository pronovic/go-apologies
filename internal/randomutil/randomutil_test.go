package randomutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomInt(t *testing.T) {
	for i := 0; i < 10000; i++ {
		r, err := RandomInt(100)
		assert.NoError(t, err)
		assert.True(t, r >= 0 && r <= 100)
	}
}

func TestRandomChoice(t *testing.T) {
	slice := []string{"one", "two", "three", "four", "five"}
	for i := 0; i < 10000; i++ {
		c, err := RandomChoice(slice)
		assert.NoError(t, err)
		assert.True(t, c == "one" || c == "two" || c == "three" || c == "four" || c == "five")
	}
}
