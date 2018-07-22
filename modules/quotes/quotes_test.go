package quotes

import (
	"github.com/lovelaced/glitzz/config"
	"testing"
)

func TestNewDoesntPanic(t *testing.T) {
	config := config.Default()
	config.QuotesDirectory = "invalid/path"
	_, err := New(nil, config)
	if err == nil {
		t.Error("error is nil")
	}
}
