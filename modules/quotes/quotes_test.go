package quotes

import (
	"github.com/lovelaced/glitzz/config"
	"testing"
)

func TestNewDoesntPanic(t *testing.T) {
	config := config.Default()
	config.Quotes.QuotesDirectory = "invalid/path"
	_, err := New(nil, config)
	if err == nil {
		t.Error("error is nil")
	}
}

func TestGetAllQuotes(t *testing.T) {
	lines, err := getAllQuotes("testdata/quotes.txt")
	if err != nil {
		t.Fatalf("error is %s", err)
	}
	if len(lines) != 3 {
		t.Errorf("returned %d lines", len(lines))
	}
}
