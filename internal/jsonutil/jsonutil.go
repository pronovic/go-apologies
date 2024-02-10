package jsonutil

import (
	"bytes"
	"encoding/json"
	"io"
)

// JSON serialization and deserialization is a bit awkward in Go.  It works reasonably well for
// public structs with public elements, but falls apart for interfaces implemented in terms of
// private structs.  For instance, take an implementation like this, which is a fairly normal way to
// implement an interface:
//
//   type Game interface {
//     PlayerCount() int
//   }
//
//   type game struct {
//     playerCount int
//   }
//
//   func NewGame(playerCount int) PlayerCount {
//     return &game { playerCount }
//   }
//
//   func (g *game) PlayerCount() int {
//     return g.playerCount
//   }
//
// The Game interface is what we expose to callers, and we want them to be able to serialize an
// instance to JSON and then round-trip back to an equivalent object.
//
// This doesn't work for a few reasons.  First, Go's JSON functionality doesn't understand
// interfaces.  Second, the game struct is private, and so are its attributes.  If we try to make
// attributes of the game struct public (i.e. changing playerCount to PlayerCount) then those names
// conflict with the existing public methods implemented for the game struct.
//
// After experimenting with a bunch of options, I think the best solution for this problem is to
// just use the awkward "X" prefix on attribute names, which avoids the conflict with public method
// names:
//
//   type game struct {
//     XplayerCount int `json:"playercount"`
//   }
//
// The second big problem pops up when you have nested interfaces (interfaces that reference other
// interfaces), like Xposition in this struct:
//
//   type pawn struct {
//     Xcolor PlayerColor `json:"playercolor"`
//     Xindex int         `json:"index"`
//     Xname string       `json:"name"`
//     Xposition Position `json:"position"`
//   }
//
// In this case, serialization generally works ok, because the JSON functionality can figure out how
// to serialize the underlying structs.  It knows that the Position interface actually points at a
// position struct, and it can invoke the standard serialization behavior for that struct.
//
// Unfortunately, JSON deserialization is not that easy.  For a simple interface like the Game shown
// above (that doesn't have any nested interfaces), you can directly decode, and the JSON
// functionality figures out what to do.  For the pawn example above (which does have nested
// interfaces), the JSON functionality has no idea how to associate the Position interface with the
// underlying position struct that implements it.
//
// Even if you implement the json.Unmarshaler interface on the position struct, it never gets
// invoked, because there is no relationship between the Position interface and the position struct.
// As far as I can tell, there is no way to register a mapping between interface and struct, or to
// provide some other sort of hint to the decoder.  As a result, there is no fully automatic way to
// to deserialize nested interfaces.
//
// One option is to make another temporary struct in terms of the underlying structs instead of the
// interfaces, like this:
//
//   type raw struct {
//     Xcolor PlayerColor `json:"playercolor"`
//     Xindex int         `json:"index"`
//     Xname string       `json:"name"`
//     Xposition position `json:"position"`
//   }
//
// You can decode the raw struct, and copy from there into a pawn struct that implements the Pawn
// interface.  This is an option if the position struct is in the same package as the pawn struct.
// However, it doesn't work if the position struct is in another package, because (by design) you
// won't have access to that struct.  It also sometimes breaks if the position struct itself
// contains nested interfaces.  Essentially, all you've done is kick the can down the road one step.
//
// After a bunch of research and failed experiments, I think the right approach is to decode JSON in
// two passes.  The first pass uses an intermediate struct that substitutes json.RawMessage for each
// nested interface:
//
//   type raw struct {
//     Xcolor PlayerColor        `json:"playercolor"`
//     Xindex int                `json:"index"`
//     Xname string              `json:"name"`
//     Xposition json.RawMessage `json:"position"`
//   }
//
// The second pass parses the json.RawMessage into an interface of the proper type.  In my code,
// this works because all of my interfaces have a constructor New<Interface>FromJSON() that accepts
// a reader and returns an interface.  We can do something like this:
//
//   var Xposition Position
//   if temp.Xposition != nil && string(temp.Xposition) != "null" {
//     Xposition, err = NewPositionFromJSON(bytes.NewReader(temp.Xposition))
//     if err != nil {
//       return nil, err
//     }
//   }
//
// Once you have the temporary Xposition variable, you can assign that into your pawn struct, and it
// will work as expected.
//
// A similar approach works for nested interfaces that are in slices or maps, although dealing with
// those is a bit more complicated, because you need to track maps or slices of json.RawMessage and
// interate through them to create maps or slices of the proper interface type.
//
// The functions below are my attempt to generalize this solution so we aren't copying around the
// same error-prone boilerplate code to every interface.  You can see how these are used by looking
// at the New<Interface>FromJSON() constructors in the model package.

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

// DecodeInterfaceJSON decodes (unmarshalls) an interface from a json.RawMessage
func DecodeInterfaceJSON[T any](raw json.RawMessage, constructor func(reader io.Reader) (T, error)) (T, error) {
	var result T
	var err error

	if raw != nil && string(raw) != "null" {
		result, err = constructor(bytes.NewReader(raw))
		if err != nil {
			return *new(T), err
		}
	}

	return result, nil
}

// DecodeMapJSON decodes (unmarshalls) a map[K]T from map[K]json.rawMessage
func DecodeMapJSON[K comparable, T any](raw map[K]json.RawMessage, constructor func(reader io.Reader) (T, error)) (map[K]T, error) {
	var result = make(map[K]T, len(raw))

	for key := range raw {
		value := raw[key]
		if value == nil || string(value) == "null" {
			var empty T
			result[key] = empty
		} else {
			element, err := constructor(bytes.NewReader(value))
			if err != nil {
				return nil, err
			}
			result[key] = element
		}
	}

	return result, nil
}

// DecodeSliceJSON decodes (unmarshalls) a []T from []json.rawMessage
func DecodeSliceJSON[T any](raw []json.RawMessage, constructor func(reader io.Reader) (T, error)) ([]T, error) {
	var result = make([]T, len(raw))

	for i := range raw {
		value := raw[i]
		if value == nil || string(value) == "null" {
			var empty T
			result[i] = empty
		} else {
			element, err := constructor(bytes.NewReader(value))
			if err != nil {
				return nil, err
			}
			result[i] = element
		}
	}

	return result, nil
}
