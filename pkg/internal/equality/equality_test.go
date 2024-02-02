package equality

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestInterface interface {
	EqualsByValue[TestInterface]
	A() int
	B() string
}

type testInterface struct {
	a int
	b string
}

func newTestInterface(a int, b string) TestInterface {
	return &testInterface{ a, b }
}

func (t *testInterface)  A() int {
	return t.a
}

func (t *testInterface)  B() string {
	return t.b
}

func (t *testInterface) Equals(other TestInterface) bool {
	return t.a == other.A() && t.b == other.B()
}

func TestByValueEquals(t *testing.T) {
	one := newTestInterface(1, "X")
	two := newTestInterface(2, "Y")
	three := newTestInterface(3, "Z")

	assert.True(t, ByValueEquals[TestInterface](nil, nil))

	assert.True(t, ByValueEquals[TestInterface](one, one))
	assert.True(t, ByValueEquals[TestInterface](two, two))
	assert.True(t, ByValueEquals[TestInterface](three, three))

	assert.False(t, ByValueEquals[TestInterface](one, nil))
	assert.False(t, ByValueEquals[TestInterface](two, nil))
	assert.False(t, ByValueEquals[TestInterface](three, nil))

	assert.False(t, ByValueEquals[TestInterface](nil, one))
	assert.False(t, ByValueEquals[TestInterface](nil, two))
	assert.False(t, ByValueEquals[TestInterface](nil, three))

	assert.False(t, ByValueEquals[TestInterface](one, two))
	assert.False(t, ByValueEquals[TestInterface](one, three))
	assert.False(t, ByValueEquals[TestInterface](two, one))
	assert.False(t, ByValueEquals[TestInterface](two, three))
	assert.False(t, ByValueEquals[TestInterface](three, one))
	assert.False(t, ByValueEquals[TestInterface](three, two))
}

func TestIntPointerEquals(t *testing.T) {
	one := 1
	two := 2
	three := 3

	assert.True(t, IntPointerEquals(nil, nil))

	assert.True(t, IntPointerEquals(&one, &one))
	assert.True(t, IntPointerEquals(&two, &two))
	assert.True(t, IntPointerEquals(&three, &three))

	assert.False(t, IntPointerEquals(&one, nil))
	assert.False(t, IntPointerEquals(&two, nil))
	assert.False(t, IntPointerEquals(&three, nil))

	assert.False(t, IntPointerEquals(nil, &one))
	assert.False(t, IntPointerEquals(nil, &two))
	assert.False(t, IntPointerEquals(nil, &three))

	assert.False(t, IntPointerEquals(&one, &two))
	assert.False(t, IntPointerEquals(&one, &three))
	assert.False(t, IntPointerEquals(&two, &one))
	assert.False(t, IntPointerEquals(&two, &three))
	assert.False(t, IntPointerEquals(&three, &one))
	assert.False(t, IntPointerEquals(&three, &two))
}

func TestStringPointerEquals(t *testing.T) {
	one := "one"
	two := "two"
	three := "three"

	assert.True(t, StringPointerEquals(nil, nil))

	assert.True(t, StringPointerEquals(&one, &one))
	assert.True(t, StringPointerEquals(&two, &two))
	assert.True(t, StringPointerEquals(&three, &three))

	assert.False(t, StringPointerEquals(&one, nil))
	assert.False(t, StringPointerEquals(&two, nil))
	assert.False(t, StringPointerEquals(&three, nil))

	assert.False(t, StringPointerEquals(nil, &one))
	assert.False(t, StringPointerEquals(nil, &two))
	assert.False(t, StringPointerEquals(nil, &three))

	assert.False(t, StringPointerEquals(&one, &two))
	assert.False(t, StringPointerEquals(&one, &three))
	assert.False(t, StringPointerEquals(&two, &one))
	assert.False(t, StringPointerEquals(&two, &three))
	assert.False(t, StringPointerEquals(&three, &one))
	assert.False(t, StringPointerEquals(&three, &two))
}