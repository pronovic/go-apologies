package identifier

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// newId() wraps NewId() to test the stubbing behavior
// note that this is not thread-safe, so it could result in instability if tests are run in parallel
func newId(stub bool) string {
	if stub {
		old := stubbed
		UseStubbedId()
		current := NewId()
		stubbed = old
		return current
	} else {
		old := stubbed
		stubbed = nil
		current := NewId()
		stubbed = old
		return current
	}
}

func TestIdImplementation(t *testing.T) {
	// This test covers the entire implementation

	stubbed := newId(true)
	assert.Equal(t, GetStubbedId(), stubbed)
	assert.Equal(t, StubbedId, stubbed)

	current := newId(false)
	assert.NotEmpty(t, current)
	assert.NotEqual(t, StubbedId, current)
}