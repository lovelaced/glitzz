package vatsim

import (
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"testing"
)

type apiMock struct {
	icao *string

	metar string
	err   error
}

func (a *apiMock) GetMetar(icao string) (string, error) {
	a.icao = &icao
	return a.metar, a.err
}

func TestMetar(t *testing.T) {
	apiMock := &apiMock{metar: "metar", err: nil}
	m, err := newVatsim(nil, config.Default(), apiMock)
	if err != nil {
		t.Fatalf("could not create the module: %s", err)
	}

	output, err := m.RunCommand(core.Command{Text: ".metar epkk"})
	if err != nil {
		t.Fatalf("error was not nil: %s", err)
	}
	if len(output) != 1 {
		t.Fatalf("invalid output length: %d", len(output))
	}
	if apiMock.icao == nil {
		t.Fatalf("icao was nil")
	}
	if *apiMock.icao != "epkk" {
		t.Fatalf("icao was invalid: %s", *apiMock.icao)
	}
}

func TestMetarNoArguments(t *testing.T) {
	m, err := newVatsim(nil, config.Default(), nil)
	if err != nil {
		t.Fatalf("could not create the module: %s", err)
	}

	output, err := m.RunCommand(core.Command{Text: ".metar"})
	if err != nil {
		t.Fatalf("error was not nil: %s", err)
	}
	if len(output) != 0 {
		t.Fatalf("invalid output length: %d", len(output))
	}
}

func TestMetarTooManyArguments(t *testing.T) {
	m, err := newVatsim(nil, config.Default(), nil)
	if err != nil {
		t.Fatalf("could not create the module: %s", err)
	}

	output, err := m.RunCommand(core.Command{Text: ".metar epkk epwa"})
	if err != nil {
		t.Fatalf("error was not nil: %s", err)
	}
	if len(output) != 0 {
		t.Fatalf("invalid output length: %d", len(output))
	}
}
