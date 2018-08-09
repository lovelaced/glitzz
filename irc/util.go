package irc

import "strings"

// IsChannelName returns true if s is a channel name.
func IsChannelName(s string) bool {
	return strings.HasPrefix(s, "#")
}
