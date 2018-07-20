package info

import (
	"github.com/lovelaced/glitzz/config"
	"strings"
	"testing"
)

func TestGit(t *testing.T) {
	p := New(nil, config.Default())
	output, err := p.RunCommand(".git")
	if err != nil {
		t.Errorf("error was not nil %s", err)
	}
	if len(output) != 1 {
		t.Errorf("invalid output length %d", len(output))
	}
	if !strings.HasPrefix(output[0], "http") {
		t.Errorf("invalid output %s", output[0])
	}
}
