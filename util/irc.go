package util

import (
	"github.com/thoj/go-ircevent"
	"strings"
)

// IsChannelName returns true if s is a channel name.
func IsChannelName(s string) bool {
	return strings.HasPrefix(s, "#")
}

// SelectReplyTarget returns the channel name if the message which triggered
// this event was sent in a channel or the sender's nickname if the message was
// sent directly to the bot.
func SelectReplyTarget(e *irc.Event) string {
	if IsChannelName(e.Arguments[0]) {
		return e.Arguments[0]
	} else {
		return e.Nick
	}
}
