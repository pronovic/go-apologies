package timestamp

import "time"

// See also: https://www.pauladamsmith.com/blog/2011/05/go_time.html

const Layout = "2006-01-02T15:04:05.000Z"
const StubbedTime = "2024-01-31T08:15:03.221Z"
var stubbed *time.Time = nil

// CurrentTime returns the current time as from time.Now().UTC()
func CurrentTime() time.Time {
	if stubbed != nil {
		return *stubbed
	} else {
		return time.Now().UTC()
	}
}

// ParseTime parses a timestamp like "2024-01-31T08:15:03.221Z" in UTC
func ParseTime(value string) (time.Time, error) {
	return time.ParseInLocation(Layout, value, time.UTC)
}

// UseStubbedTime sets the current time to "2024-01-31T08:15:03.221Z" (as UTC), for use in unit tests
func UseStubbedTime() {
	parsed, _ := ParseTime(StubbedTime)
	stubbed = &parsed
}

// GetStubbedTime returns the stubbed timestamp, for use in unit tests
// Note that this is the expected value, which may or may not actually be in use
func GetStubbedTime() time.Time {
	parsed, _ := ParseTime(StubbedTime)
	return parsed
}