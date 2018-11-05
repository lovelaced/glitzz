package decide

import (
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"testing"
)

func TestDecideNoParameters(t *testing.T) {
	p, err := New(nil, config.Default())
	if err != nil {
		t.Fatalf("error creating module %s", err)
	}

	output, err := p.RunCommand(core.Command{Text: ".decide"})
	if err != nil {
		t.Errorf("error was not nil %s", err)
	}
	if len(output) != 0 {
		t.Errorf("output length %d", len(output))
	}
}

func TestDecideMissingOr(t *testing.T) {
	p, err := New(nil, config.Default())
	if err != nil {
		t.Fatalf("error creating module %s", err)
	}

	output, err := p.RunCommand(core.Command{Text: ".decide another thing"})
	if err != nil {
		t.Errorf("error was not nil %s", err)
	}
	if len(output) != 0 {
		t.Errorf("invalid output length %d", len(output))
	}
}

func TestDecide(t *testing.T) {
	p, err := New(nil, config.Default())
	if err != nil {
		t.Fatalf("error creating module %s", err)
	}

	output, err := p.RunCommand(core.Command{Text: ".decide thing, another thing or something else", Nick: "nick"})
	if err != nil {
		t.Errorf("error was not nil %s", err)
	}
	if len(output) != 1 {
		t.Errorf("invalid output length %d", len(output))
	}
}
