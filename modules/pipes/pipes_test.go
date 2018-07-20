package pipes

import (
	"github.com/lovelaced/glitzz/config"
	"testing"
)

func TestUpper(t *testing.T) {
	p := New(nil, config.Default())
	output, err := p.RunCommand(".upper text TEXT text")
	if err != nil {
		t.Errorf("error was not nil %s", err)
	}
	if len(output) != 1 {
		t.Errorf("invalid output length %d", len(output))
	}
	if output[0] != "TEXT TEXT TEXT" {
		t.Errorf("invalid output %s", output[0])
	}
}

func TestLower(t *testing.T) {
	p := New(nil, config.Default())
	output, err := p.RunCommand(".lower text TEXT text")
	if err != nil {
		t.Errorf("error was not nil %s", err)
	}
	if len(output) != 1 {
		t.Errorf("invalid output length %d", len(output))
	}
	if output[0] != "text text text" {
		t.Errorf("invalid output %s", output[0])
	}
}

func TestEcho(t *testing.T) {
	p := New(nil, config.Default())
	output, err := p.RunCommand(".echo text TEXT text")
	if err != nil {
		t.Errorf("error was not nil %s", err)
	}
	if len(output) != 1 {
		t.Errorf("invalid output length %d", len(output))
	}
	if output[0] != "text TEXT text" {
		t.Errorf("invalid output %s", output[0])
	}
}
