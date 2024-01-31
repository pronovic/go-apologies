package timestamp

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// currentTime() wraps CurrentTime() to test the stubbing behavior
// note that this is not thread-safe, so it could result in instability if tests are run in parallel
func currentTime(stub bool) time.Time {
	if stub {
		old := stubbed
		UseStubbedTime()
		current := CurrentTime()
		stubbed = old
		return current
	} else {
		old := stubbed
		stubbed = nil
		current := CurrentTime()
		stubbed = old
		return current
	}
}

func TestTimeImplementation(t *testing.T) {
	// This test covers the entire implementation

	now := time.Now()

	stubbed := currentTime(true)
	assert.Equal(t, GetStubbedTime(), stubbed)
	assert.Equal(t, StubbedTime, stubbed.Format(Layout))

	current := currentTime(false)
	assert.True(t, current.After(now) || current.Equal(now))
	_, offset := current.Zone()
	assert.Equal(t, 0, offset)  // zero offset means UTC
}