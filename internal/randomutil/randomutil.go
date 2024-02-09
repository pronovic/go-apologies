package randomutil

import (
	"crypto/rand"
	"math/big"
)

// RandomInt returns a random integer between zero and a max value (inclusive)
func RandomInt(max int) (int, error) {
	value, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}

	return int(value.Int64()), nil
}

// RandomChoice returns a random choice from a slice
func RandomChoice[T any](slice []T) (T, error) {
	index, err := RandomInt(len(slice))
	if err != nil {
		return *new(T), err
	}

	return slice[index], nil
}
