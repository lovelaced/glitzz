package util

import (
	"testing"
)

func TestGetRandomArrayElementEmpty(t *testing.T) {
	_, err := GetRandomArrayElement([]string{})
	if err == nil {
		t.Error("error was nil")
	}
}

func TestGetRandomArrayElement(t *testing.T) {
	_, err := GetRandomArrayElement([]string{"a", "b"})
	if err != nil {
		t.Errorf("error was %s", err)
	}
}
