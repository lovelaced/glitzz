package pornhub

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	pornhubURL = "https://www.pornhub.com"
	userAgent  = "Mozilla/5.0 (X11; Linux x86_64; rv:63.0) Gecko/20100101 Firefox/63.0"
)

type pornhub struct {
	url string
}

func GetRandPage() (err error) {
	client := &http.Client{Timeout: 10 * time.Second}

	request, err := http.NewRequest("GET", pornhubURL+"/random", nil)
	if err != nil {
		return
	}
	request.Header.Set("Accept-Language", "en-US,en;q=0.5")
	request.Header.Set("User-Agent", userAgent)

	response, err := client.Do(request)
	if err != nil {
		return
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Something went wrong: %d code", response.StatusCode)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	defer response.Body.Close()

	println(string(body))
	return
}
