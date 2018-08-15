package seen

import (
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"github.com/thoj/go-ircevent"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSeenNotFound(t *testing.T) {
	// setup
	tmpDirName, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Could not create a temporary directory: %s", err)
	}
	defer os.Remove(tmpDirName)

	conf := config.Default()
	conf.Seen.SeenFile = filepath.Join(tmpDirName, "seen_data.json")

	m, err := New(nil, conf)
	if err != nil {
		t.Fatalf("Could not create a module: %s", err)
	}

	// test
	output, err := m.RunCommand(core.Command{Text: ".seen nick", Target: "#channel"})
	if err != nil {
		t.Errorf("error was not nil %s", err)
	}
	if len(output) != 1 {
		t.Errorf("invalid output length %d", len(output))
	}
	if !strings.Contains(output[0], "have never seen") {
		t.Errorf("invalid output %s", output[0])
	}
}

func TestSeenFound(t *testing.T) {
	// setup
	tmpDirName, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Could not create a temporary directory: %s", err)
	}
	defer os.Remove(tmpDirName)

	conf := config.Default()
	conf.Seen.SeenFile = filepath.Join(tmpDirName, "seen_data.json")

	m, err := New(nil, conf)
	if err != nil {
		t.Fatalf("Could not create a module: %s", err)
	}

	// test
	e := &irc.Event{Nick: "nick", Code: "PRIVMSG", Arguments: []string{"#channel", "message"}}
	m.HandleEvent(e)

	output, err := m.RunCommand(core.Command{Text: ".seen nick", Target: "#channel"})
	if err != nil {
		t.Errorf("error was not nil %s", err)
	}
	if len(output) != 1 {
		t.Errorf("invalid output length %d", len(output))
	}
	if !strings.Contains(output[0], "was last seen") {
		t.Errorf("invalid output %s", output[0])
	}

}
