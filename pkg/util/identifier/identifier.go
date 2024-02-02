package identifier

import (
	"github.com/google/uuid"
)

const StubbedId = "id"
var stubbed *string = nil

// NewId returns a new random identifier
func NewId() string {
	if stubbed != nil {
		return *stubbed
	} else {
		return uuid.New().String()
	}
}

// UseStubbedId sets things up so that every call to NewId() returns a stubbed identifier
func UseStubbedId() {
	id := StubbedId
	stubbed = &id
}

// GetStubbedId returns the stubbed identifier, for use in unit tests
// Note that this is the expected value, which may or may not actually be in use
func GetStubbedId() string {
	return StubbedId
}