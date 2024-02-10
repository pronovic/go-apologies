package enum

import (
	"errors"
)

// GoLang does not have a concept of enumerations, like you would find in Java or Python.
//
// The canonical example suggests making a type (i.e. `type MyEnum string` or `type MyEnum int`) and
// then creating constants of type MyEnum.  This (sort of) gets you type safety in methods that need
// to accept an enumeration, but doesn't prevent callers from making up other values of MyEnum that
// you didn't expect.  Also, there's no way to actually enumerate your enum: you can't see what
// values are part of it.  To make matters worse, those names also share a namespace with all
// other constants and variables defined in the same scope.
//
// There's no general way to fix the shared namespace problem, other than prefixing each enumeration
// constant with the name of the enumeration (yuck).  However, we can do something about the other
// problems.
//
// This is based conceptually on the pattern at https://github.com/nikolaydubina/go-enum-example,
// with further discussion found here: https://news.ycombinator.com/item?id=37703175.  There are
// lots of other StackOverflow articles and Reddit discussions on the topic, but this is what I
// ended up with.
//
// Individual enumerations are defined as simple structs with a private member, and they all
// implement the Enum interface with its Value() method:
//
//    type GameMode struct { value string }
//    func (e GameMode) Value() string { return e.value }
//
// Since the value element of the struct is private, no code outside your package will be able to
// create other variables of the type MyEnum.  Implementing the Enum interface makes it possible to
// share common functionality implemented below, such as tracking the values of an enum.
//
// Next, we define values of the enumeration:
//
//    var GameModes = enum.NewValues[GameMode](StandardMode, AdultMode)
//    var StandardMode = GameMode{"Standard"}
//    var AdultMode = GameMode{"Adult"}
//
// The convention I am following is that the enumeration itself (GameMode) is singular, and the
// variable that holds its values is plural (GameModes). Unfortunately, these all need to be var
// instead of const, because you can't assign a struct to a const.  In this case, I'm assuming there
// won't be a conflict with names like Standard and Adult, but if there were I could prefix the
// enumeration name (i.e. GameModeStandard) at the cost of some readability.
//
// Finally, if you want to use the enumeration in JSON, you should implement two additional methods:
//
//   func (e GameMode) MarshalText() (text []byte, err error) { return enum.Marshal(e) }
//   func (e *GameMode) UnmarshalText(text []byte) error { return enum.Unmarshal(e, text, GameModes) }
//
// Note that the unmarshal interfaces always expect a pointer receiver, so you can change the receiver
// to point at the resulting unmarshalled value.  This is a bit awkward, because the receiver for
// Value() must be a non-pointer.  There's no good way around this, so I ignore the Goland warning.
//
// Overall, this implementation isn't nearly as full-featured as Enum in Python, but it gets me enough
// type-safety for my purposes, plus most of the functionality I rely on regularly.

// Enum is a member of an enumeration
type Enum interface {
	Value() string
}

// Values is a fixed collection of enumeration values
type Values[T Enum] interface {

	// Members returns all members of an enumeration
	Members() []T

	// MemberValues returns all values that make up an enumeration
	MemberValues() []string

	// MemberOf checks whether a value is a member of the enumeration
	MemberOf(value string) bool

	// GetMember returns the member with the provided value, or an error if no such member exists
	GetMember(value string) (T, error)
}

type values[T Enum] struct {
	members []T
}

// NewValues constructs a new Values collection with a set of members
func NewValues[T Enum](members ...T) Values[T] {
	return &values[T]{append(make([]T, 0, len(members)), members...)}
}

func (v *values[T]) Members() []T {
	return append(make([]T, 0, len(v.members)), v.members...)
}

func (v *values[T]) MemberValues() []string {
	members := make([]string, 0, len(v.members))

	for _, m := range v.members {
		members = append(members, m.Value())
	}

	return members
}

func (v *values[T]) MemberOf(value string) bool {
	for _, m := range v.members {
		if m.Value() == value {
			return true
		}
	}

	return false
}

func (v *values[T]) GetMember(value string) (T, error) {
	for _, m := range v.members {
		if m.Value() == value {
			return m, nil
		}
	}

	return *new(T), errors.New("member not found")
}

// Marshal marshals the enum value to text, useful when implementing MarshalText or MarshalJSON
func Marshal(e Enum) (text []byte, err error) {
	return []byte(e.Value()), nil
}

// Unmarshal unmarshals text into an enum value, useful when implementing UnmarshalText or UnmarshalJSON
func Unmarshal[T Enum](e *T, text []byte, values Values[T]) error {
	value, err := values.GetMember(string(text[:]))
	if err != nil {
		return err
	}

	*e = value
	return nil
}
