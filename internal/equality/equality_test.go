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

func TestSliceByValueEquals(t *testing.T) {
	one := newTestInterface(1, "X")
	two := newTestInterface(2, "Y")
	three := newTestInterface(3, "Z")

	slice1 := []TestInterface { one }
	slice2 := []TestInterface { one }
	slice3 := []TestInterface { two, three }
	slice4 := []TestInterface { two, one }
	slice5 := []TestInterface { two, one }

	assert.True(t, SliceByValueEquals[TestInterface](nil, nil))
	assert.True(t, SliceByValueEquals[TestInterface](slice1, slice1))
	assert.True(t, SliceByValueEquals[TestInterface](slice1, slice2))
	assert.True(t, SliceByValueEquals[TestInterface](slice2, slice2))
	assert.True(t, SliceByValueEquals[TestInterface](slice2, slice1))
	assert.True(t, SliceByValueEquals[TestInterface](slice3, slice3))
	assert.True(t, SliceByValueEquals[TestInterface](slice4, slice4))
	assert.True(t, SliceByValueEquals[TestInterface](slice4, slice5))
	assert.True(t, SliceByValueEquals[TestInterface](slice5, slice5))
	assert.True(t, SliceByValueEquals[TestInterface](slice5, slice4))

	assert.False(t, SliceByValueEquals[TestInterface](nil, slice1))
	assert.False(t, SliceByValueEquals[TestInterface](nil, slice2))
	assert.False(t, SliceByValueEquals[TestInterface](nil, slice3))
	assert.False(t, SliceByValueEquals[TestInterface](nil, slice4))
	assert.False(t, SliceByValueEquals[TestInterface](nil, slice5))

	assert.False(t, SliceByValueEquals[TestInterface](slice1, nil))
	assert.False(t, SliceByValueEquals[TestInterface](slice1, slice3))
	assert.False(t, SliceByValueEquals[TestInterface](slice1, slice4))
	assert.False(t, SliceByValueEquals[TestInterface](slice1, slice5))

	assert.False(t, SliceByValueEquals[TestInterface](slice3, nil))
	assert.False(t, SliceByValueEquals[TestInterface](slice3, slice1))
	assert.False(t, SliceByValueEquals[TestInterface](slice3, slice2))
	assert.False(t, SliceByValueEquals[TestInterface](slice3, slice4))
	assert.False(t, SliceByValueEquals[TestInterface](slice3, slice5))

	assert.False(t, SliceByValueEquals[TestInterface](slice4, nil))
	assert.False(t, SliceByValueEquals[TestInterface](slice4, slice1))
	assert.False(t, SliceByValueEquals[TestInterface](slice4, slice2))
	assert.False(t, SliceByValueEquals[TestInterface](slice4, slice3))

	assert.False(t, SliceByValueEquals[TestInterface](slice5, nil))
	assert.False(t, SliceByValueEquals[TestInterface](slice5, slice1))
	assert.False(t, SliceByValueEquals[TestInterface](slice5, slice2))
	assert.False(t, SliceByValueEquals[TestInterface](slice5, slice3))
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