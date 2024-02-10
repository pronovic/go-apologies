package enum

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Sample struct{ value string }

func (e Sample) Value() string                         { return e.value }
func (e Sample) MarshalText() (text []byte, err error) { return Marshal(e) }
func (e *Sample) UnmarshalText(text []byte) error      { return Unmarshal(e, text, SampleValues) }

var Value1 = Sample{"Value1"}
var Value2 = Sample{"Value2"}
var SampleValues = NewValues[Sample](Value1, Value2)

func TestMembers(t *testing.T) {
	assert.Equal(t, []Sample{Value1, Value2}, SampleValues.Members())
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
	var member Sample

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

func TestJsonMarshalUnmarshal(t *testing.T) {
	value := Value2
	marshalled, err := json.Marshal(value)
	assert.Nil(t, err)
	var unmarshalled Sample
	err = json.Unmarshal(marshalled, &unmarshalled)
	assert.Nil(t, err)
	assert.True(t, value == unmarshalled)
}
