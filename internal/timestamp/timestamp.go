package timestamp

import (
	"time"
)

// See also: https://www.pauladamsmith.com/blog/2011/05/go_time.html

const Layout = "2006-01-02T15:04:05.000Z"

type Factory interface {
	CurrentTime() time.Time
}

type factory struct {
}

func NewFactory() Factory {
	return &factory{}
}

// CurrentTime returns the current time as from time.Now().UTC()
func (f *factory) CurrentTime() time.Time {
	return time.Now().UTC()
}

// ParseTime parses a timestamp like "2024-01-31T08:15:03.221Z" in UTC
func ParseTime(value string) (time.Time, error) {
	return time.ParseInLocation(Layout, value, time.UTC)
}

// FormatTime formats a timestamp like "2024-01-31T08:15:03.221Z" in UTC
func FormatTime(value time.Time) string {
	return value.UTC().Format(Layout)
}

// Marshal marshals a time value to text, useful when implementing MarshalText or MarshalJSON
func Marshal(t time.Time) (text []byte, err error) {
	return []byte(FormatTime(t)), nil
}

// Unmarshal unmarshals text into a time value, useful when implementing UnmarshalText or UnmarshalJSON
func Unmarshal(t *time.Time, text []byte) error {
	value, err := ParseTime(string(text[:]))
	if err != nil {
		return err
	}

	*t = value
	return nil
}