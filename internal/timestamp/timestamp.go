package timestamp

import "time"

// See also: https://www.pauladamsmith.com/blog/2011/05/go_time.html

const Layout = "2006-01-02T15:04:05.000Z"

type Timestamp time.Time

type Factory interface {
	CurrentTime() Timestamp
}

type factory struct{}

func NewFactory() Factory {
	return &factory{}
}

// CurrentTime returns the current time as from time.Now().UTC()
func (f *factory) CurrentTime() Timestamp {
	// The format/parse step normalizes the date, removing the nanosecond precision that we
	// don't use in the standard layout (i.e. 561093000ns vs. 561000000ns).  Without this,
	// serialization back and forth to text is lossy and we get spurious differences.
	t := Timestamp(time.Now().UTC())
	n, _ := Parse((&t).Format())
	return n
}

// Parse parses a timestamp like "2024-01-31T08:15:03.221Z" in UTC
func Parse(value string) (Timestamp, error) {
	t, err := time.ParseInLocation(Layout, value, time.UTC)
	return Timestamp(t), err
}

// Format formats a timestamp like "2024-01-31T08:15:03.221Z" in UTC
func (t *Timestamp) Format() string {
	return t.AsTime().UTC().Format(Layout)
}

// AsTime converts the timestamp to a standard time.Time
func (t *Timestamp) AsTime() time.Time {
	return time.Time(*t)
}

// MarshalText marshals a time value to JSON
func (t Timestamp) MarshalText() (text []byte, err error) {
	return []byte(t.Format()), nil
}

// UnmarshalText unmarshals text into a timestamp value
func (t *Timestamp) UnmarshalText(text []byte) error {
	content := string(text[:])

	value, err := Parse(content)
	if err != nil {
		return err
	}

	*t = value
	return nil
}
