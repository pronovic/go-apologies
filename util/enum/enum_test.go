package enum

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Sample struct{ value string }

func (e Sample) Value() string { return e.value }

var Value1 = Sample{"Value1"}
var Value2 = Sample{"Value2"}
var SampleValues = Values[Sample](Value1, Value2)

func TestMembers(t *testing.T) {
	assert.Equal(t, []Enum{Value1, Value2}, SampleValues.members)
}

func TestMemberValues(t *testing.T) {
	assert.Equal(t, []string{"Value1", "Value2"}, SampleValues.MemberValues())
}

func TestMemberOf(t *testing.T) {
	assert.Equal(t, false, SampleValues.MemberOf(""))
	assert.Equal(t, false, SampleValues.MemberOf("bogus"))
	assert.Equal(t, true, SampleValues.MemberOf("Value1"))
	assert.Equal(t, true, SampleValues.MemberOf("Value2"))
}

func TestGetMember(t *testing.T) {
	var err error
	var member Enum

	_, err = SampleValues.GetMember("")
	assert.EqualError(t, err, "member not found")

	_, err = SampleValues.GetMember("bogus")
	assert.EqualError(t, err, "member not found")

	member, err = SampleValues.GetMember("Value1")
	assert.Nil(t, err)
	assert.Equal(t, Value1, member)

	member, err = SampleValues.GetMember("Value2")
	assert.Nil(t, err)
	assert.Equal(t, Value2, member)
}
