package equality

import (
	"fmt"
	"testing"
	"time"
)

// This is based on testify's assert.Equal() unit test.
// Original source: https://github.com/stretchr/testify/blob/c3b0c9b4f50cd8320ccdcdfd3ffd6afc5b109c4a/assert/assertions_test.go#L104
func TestEqualByValue(t *testing.T) {
	cases := []struct {
		expected interface{}
		actual   interface{}
		result   bool
	}{
		// cases that are expected to be equal
		{"Hello World", "Hello World", true},
		{123, 123, true},
		{123.5, 123.5, true},
		{[]byte("Hello World"), []byte("Hello World"), true},
		{nil, nil, true},

		// cases that are expected not to be equal
		{map[int]int{5: 10}, map[int]int{10: 20}, false},
		{'x', "x", false},
		{"x", 'x', false},
		{0, 0.1, false},
		{0.1, 0, false},
		{time.Now, time.Now, false},
		{func() {}, func() {}, false},
		{uint32(10), int32(10), false},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("EqualByValue(%#v, %#v)", c.expected, c.actual), func(t *testing.T) {
			res := EqualByValue(c.expected, c.actual)

			if res != c.result {
				t.Errorf("EqualByValue(%#v, %#v) should return %#v", c.expected, c.actual, c.result)
			}
		})
	}
}