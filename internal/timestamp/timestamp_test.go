package timestamp

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFactoryCurrentTime(t *testing.T) {
	now := time.Now()
	current := NewFactory().CurrentTime()
	assert.True(t, current.After(now) || current.Equal(now))
	_, offset := current.Zone()
	assert.Equal(t, 0, offset)  // zero offset means UTC
}

func TestParseTime(t *testing.T) {
	input := "2024-01-31T08:15:03.221Z"
	parsed, err := ParseTime(input)
	assert.Nil(t, err)
	assert.Equal(t, time.Date(2024, time.January, 31, 8, 15, 3, 221000000, time.UTC), parsed)
	_, offset := parsed.Zone()
	assert.Equal(t, 0, offset)  // zero offset means UTC
}