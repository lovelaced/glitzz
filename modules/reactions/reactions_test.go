package reactions

import (
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"testing"
)

func TestCuteNoParameters(t *testing.T) {
	p, err := New(nil, config.Default())
	if err != nil {
		t.Fatalf("error creating module %s", err)
	}

	output, err := p.RunCommand(core.Command{Text: ".cute", Nick: "nick"})
	if err != nil {
		t.Errorf("error was not nil %s", err)
	}
	if len(output) != 1 {
		t.Errorf("invalid output length %d", len(output))
	}
}

func TestCute(t *testing.T) {
	p, err := New(nil, config.Default())
	if err != nil {
		t.Fatalf("error creating module %s", err)
	}

	output, err := p.RunCommand(core.Command{Text: ".cute param1 param2", Nick: "nick"})
	if err != nil {
		t.Errorf("error was not nil %s", err)
	}
	if len(output) != 1 {
		t.Errorf("invalid output length %d", len(output))
	}
}

func TestMagicNoParameters(t *testing.T) {
	p, err := New(nil, config.Default())
	if err != nil {
		t.Fatalf("error creating module %s", err)
	}

	_, err = p.RunCommand(core.Command{Text: ".magic", Nick: "nick"})
	if err == nil {
		t.Error("error was nil")
	}
}

func TestMagic(t *testing.T) {
	p, err := New(nil, config.Default())
	if err != nil {
		t.Fatalf("error creating module %s", err)
	}

	output, err := p.RunCommand(core.Command{Text: ".magic param1 param2", Nick: "nick"})
	if err != nil {
		t.Errorf("error was not nil %s", err)
	}
	if len(output) != 1 {
		t.Errorf("invalid output length %d", len(output))
	}
}

func TestStumpNoParameters(t *testing.T) {
	p, err := New(nil, config.Default())
	if err != nil {
		t.Fatalf("error creating module %s", err)
	}

	_, err = p.RunCommand(core.Command{Text: ".stump", Nick: "nick"})
	if err == nil {
		t.Error("error was nil")
	}
}

func TestStump(t *testing.T) {
	p, err := New(nil, config.Default())
	if err != nil {
		t.Fatalf("error creating module %s", err)
	}

	output, err := p.RunCommand(core.Command{Text: ".stump param1 param2", Nick: "nick"})
	if err != nil {
		t.Errorf("error was not nil %s", err)
	}
	if len(output) != 1 {
		t.Errorf("invalid output length %d", len(output))
	}
}

func TestSpurd(t *testing.T) {
	p, err := New(nil, config.Default())
	if err != nil {
		t.Fatalf("error creating module %s", err)
	}

	output, err := p.RunCommand(core.Command{Text: ".spurd lorem ipsum", Nick: "nick"})
	if err != nil {
		t.Errorf("error was not nil %s", err)
	}
	if len(output) != 1 {
		t.Errorf("invalid output length %d", len(output))
	}
}

func TestSpurdNoParameters(t *testing.T) {
	p, err := New(nil, config.Default())
	if err != nil {
		t.Fatalf("error creating module %s", err)
	}

	output, err := p.RunCommand(core.Command{Text: ".spurd", Nick: "nick"})
	if err != nil {
		t.Errorf("error was not nil %s", err)
	}
	if len(output) != 0 {
		t.Errorf("invalid output length %d", len(output))
	}
}

func TestReplaceAndPreserveCaseLower(t *testing.T) {
	output := replaceAndPreserveCase("lorem", "or", "aa")
	if output != "laaem" {
		t.Errorf("output was: %s", output)
	}
}

func TestReplaceAndPreserveCaseMultiple(t *testing.T) {
	output := replaceAndPreserveCase("loremor", "or", "aa")
	if output != "laaemaa" {
		t.Errorf("output was: %s", output)
	}
}

func TestReplaceAndPreserveCaseUpper(t *testing.T) {
	output := replaceAndPreserveCase("LOREM", "or", "aa")
	if output != "LAAEM" {
		t.Errorf("output was: %s", output)
	}
}

func TestReplaceAndPreserveCaseMixed(t *testing.T) {
	output := replaceAndPreserveCase("LoReM", "or", "aa")
	if output != "LaAeM" {
		t.Errorf("output was: %s", output)
	}
}

func TestReplaceAndPreserveCaseLongerUnicode(t *testing.T) {
	output := replaceAndPreserveCase("LśĆeM", "ść", "aaaa")
	if output != "LaAaaeM" {
		t.Errorf("output was: %s", output)
	}
}

func TestReplaceAndPreserveCaseLonger(t *testing.T) {
	output := replaceAndPreserveCase("LoReM", "or", "aaaa")
	if output != "LaAaaeM" {
		t.Errorf("output was: %s", output)
	}
}

type spurdReplaceTest struct {
	Input  string
	Output string
}

var spurdReplaceTests = []spurdReplaceTest{
	{
		Input:  "SimpleFold iterates over Unicode code points equivalent under the Unicode-defined simple case folding.",
		Output: "SimpleFold iderates over Unigode gode boidnz equibaledn under de Unigode-defined simple gase foldign.",
	},
	{
		Input:  "SIMPLEFOLD ITERATES OVER UNICODE CODE POINTS EQUIVALENT UNDER THE UNICODE-DEFINED SIMPLE CASE FOLDING.",
		Output: "SIMPLEFOLD IDERATES OVER UNIGODE GODE BOIDNZ EQUIBALEDN UNDER DE UNIGODE-DEFINED SIMPLE GASE FOLDIGN.",
	},
}

func TestSpurdReplace(t *testing.T) {
	for _, test := range spurdReplaceTests {
		output := spurdReplace(test.Input)
		if output != test.Output {
			t.Logf("Expected: %+v", test.Output)
			t.Logf("Actual:   %+v", output)
			t.Error("outputs differ")
		}
	}
}

func TestSpurdReplacements(t *testing.T) {
	for _, elem := range spurdReplacements {
		if len(elem) != 2 {
			t.Errorf("replacement should have two elements: %+v", elem)
		}
	}
}
