package pornmd

import (
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

type pornmdResults []struct {
	Segment string `json:"segment"`
	Keyword string `json:"keyword"`
}

// ReturnRandSearch returns a random search item
func ReturnRandSearch() ([]string, error) {
	var genre, result string
	var mdLiveTerms pornmdResults
	client := &http.Client{Timeout: 10 * time.Second}

	response, err := client.Get(pornmdLiveTermsURL)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return []string{fmt.Sprintf("Something went wrong: %d code", response.StatusCode)}, nil
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	err = json.Unmarshal(body, &mdLiveTerms)
	if err != nil {
		return nil, err
	}
	if len(mdLiveTerms) == 0 {
		return nil, errors.New("No results, wtf?")
	}

	index := rand.Intn(len(mdLiveTerms))
	switch mdLiveTerms[index].Segment {
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
		mdLiveTerms[index].Keyword,
		genre)

	return []string{result}, nil
}
