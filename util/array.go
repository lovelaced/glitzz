package util

import (
	"errors"
	"math/rand"
)

func GetRandomArrayElement(array []string) (string, error) {
	if len(array) == 0 {
		return "", errors.New("array length is zero")
	}
	return array[rand.Intn(len(array))], nil
}
