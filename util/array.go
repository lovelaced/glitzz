package util

import (
	"errors"
	"math/rand"
)

// GetRandomArrayElement selects a random element from the array. If the array
// is empty an error is returned.
func GetRandomArrayElement(array []string) (string, error) {
	if len(array) == 0 {
		return "", errors.New("array length is zero")
	}
	return array[rand.Intn(len(array))], nil
}
