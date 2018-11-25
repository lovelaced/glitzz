package porn

import (
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"

	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

const (
	pornmdLiveTermsURL = "https://www.pornmd.com/getliveterms"
)

// New registers the stackexchange module
func New(sender core.Sender, conf config.Config) (core.Module, error) {
	rv := &porn{
		Base: core.NewBase("porn", sender, conf),
	}

	rv.AddCommand("pornmd", rv.pornmd)
	return rv, nil
}

type results []struct {
	Segment string `json:"segment"`
	Keyword string `json:"keyword"`
}

type porn struct {
	core.Base
	liveterms results
}

func (p *porn) pornmd(arguments core.CommandArguments) ([]string, error) {
	var genre, result string
	client := http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(pornmdLiveTermsURL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return []string{fmt.Sprintf("Something went wrong: %d code", resp.StatusCode)}, nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	err = json.Unmarshal(body, &p.liveterms)
	if err != nil {
		return nil, err
	}
	if len(p.liveterms) == 0 {
		return nil, errors.New("No results, wtf?")
	}

	index := rand.Intn(len(p.liveterms))
	switch p.liveterms[index].Segment {
	case "g":
		genre = "Gay"
	case "t":
		genre = "Trap"
	case "s":
		genre = "Straight"
	default:
		genre = "(NULL)"
	}
	result = fmt.Sprintf("%s | %s",
		p.liveterms[index].Keyword,
		genre)

	return []string{result}, nil
}
