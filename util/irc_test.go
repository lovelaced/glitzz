package util

import (
	"testing"
)

func TestIsChannelName(t *testing.T) {
	if !IsChannelName("#channel") {
		t.Errorf("#channel should be a channel name")
	}

	if IsChannelName("nick") {
		t.Errorf("nick should not be a channel name")
	}
}
