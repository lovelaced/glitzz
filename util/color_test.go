package util

import (
	"strings"
	"testing"
)

func TestGreentext(t *testing.T) {
	input := "input"
	output := Greentext(input)
	if !strings.Contains(output, input) {
		t.Errorf("invalid output: %s", output)
	}
	if output == input {
		t.Errorf("noting happened, output: %x, input: %x", output, input)
	}
}

func TestNormaltext(t *testing.T) {
	input := "input"
	output := Normaltext(input)
	if !strings.Contains(output, input) {
		t.Errorf("invalid output: %s", output)
	}
	if output == input {
		t.Errorf("noting happened, output: %x, input: %x", output, input)
	}
}
