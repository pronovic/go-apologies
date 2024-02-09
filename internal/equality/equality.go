package equality

import (
	"bytes"
	"reflect"
)

// EqualByValue determines if two objects are considered equal.
// This is based on testify's assert.Equal() implementation, which is general and seems to be reliable.
// Original source: https://github.com/stretchr/testify/blob/c3b0c9b4f50cd8320ccdcdfd3ffd6afc5b109c4a/assert/assertions.go#L58
func EqualByValue(left any, right any) bool {
	if left == nil || right == nil {
		return left == right
	}

	l, ok := left.([]byte)
	if !ok {
		return reflect.DeepEqual(left, right)
	}

	r, ok := right.([]byte)
	if !ok {
		return false
	}

	if l == nil || r == nil {
		return l == nil && r == nil
	}

	return bytes.Equal(l, r)
}