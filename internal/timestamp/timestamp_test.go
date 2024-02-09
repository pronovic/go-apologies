package timestamp

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFactoryCurrentTime(t *testing.T) {
	now := time.Now()
	current := NewFactory().CurrentTime()
	assert.True(t, current.AsTime().After(now) || current.AsTime().Equal(now))
	_, offset := current.AsTime().Zone()
	assert.Equal(t, 0, offset)  // zero offset means UTC
}

func TestParseFormatCurrent(t *testing.T) {
	// it's important that we can reliably round-trip without spurious differences
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
	assert.Equal(t, 0, offset)  // zero offset means UTC

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