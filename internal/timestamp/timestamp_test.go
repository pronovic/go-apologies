package timestamp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFactoryCurrentTime(t *testing.T) {
	// The check vs. now needs to be for a millisecond before the current time, because
	// our timestamp format stores only millisecond precision.  If the current time has
	// 561093000ns, then we'll end up with 561000000ns.  That means we can't reliably
	// test against After(now) unless we subtract one 1ms from now (getting us 560000000ns).
	now := time.Now().Add(-1 * time.Millisecond)
	current := NewFactory().CurrentTime()
	assert.True(t, current.AsTime().After(now))
	_, offset := current.AsTime().Zone()
	assert.Equal(t, 0, offset) // zero offset means UTC
}

func TestParseFormatCurrent(t *testing.T) {
	// It's important that we can reliably round-trip without spurious differences.
	// The string format only has millisecond precision, so nanoseconds are removed
	// from the returned timestamp.  This test confirms that round-trip works.
	current := NewFactory().CurrentTime()
	parsed, _ := Parse(current.Format())
	assert.Equal(t, current, parsed)
}

func TestParseFormat(t *testing.T) {
	input := "2024-01-31T08:15:03.221Z"

	parsed, err := Parse(input)
	assert.Nil(t, err)
	assert.Equal(t, time.Date(2024, time.January, 31, 8, 15, 3, 221000000, time.UTC), parsed.AsTime())
	_, offset := parsed.AsTime().Zone()
	assert.Equal(t, 0, offset) // zero offset means UTC

	formatted := parsed.Format()
	assert.Equal(t, input, formatted)
}

func TestMarshalUnmarshal(t *testing.T) {
	input := "2024-01-31T08:15:03.221Z"
	parsed, _ := Parse(input)

	marshalled, err := parsed.MarshalText()
	assert.Nil(t, err)
	assert.Equal(t, input, string(marshalled))
	var unmarshalled Timestamp
	err = unmarshalled.UnmarshalText(marshalled)
	assert.Nil(t, err)
	assert.Equal(t, parsed, unmarshalled)
}
