package identifier

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFactoryRandomId(t *testing.T) {
	one := NewFactory().RandomId()
	two := NewFactory().RandomId()
	assert.NotEmpty(t, one)
	assert.NotEmpty(t, two)
	assert.NotEqual(t, one, two)
}