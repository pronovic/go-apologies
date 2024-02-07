package identifier

import (
	"github.com/google/uuid"
)

type Factory interface {
	RandomId() string
}

type factory struct {
}

func NewFactory() Factory {
	return &factory{}
}

// RandomId returns a new random identifier
func (f *factory) RandomId() string {
	return uuid.New().String()
}