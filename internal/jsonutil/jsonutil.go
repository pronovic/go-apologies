package jsonutil

import (
	"encoding/json"
	"io"
)

// DecodeSimpleJSON decodes (unmarshalls) JSON for a struct containing only simple data types.
// This works fine if your struct has only int, string, bool, enumerations, etc.  If your struct
// contains nested interfaces, it does *not* work and you need a more complicated implementation.
func DecodeSimpleJSON[T any](reader io.Reader) (*T, error) {
	var obj T

	err := json.NewDecoder(reader).Decode(&obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}
