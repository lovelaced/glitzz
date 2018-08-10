package vatsim

import (
	"errors"
	"fmt"
	"github.com/lovelaced/glitzz/logging"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type vatsimApi interface {
	GetMetar(icao string) (string, error)
}

const metarApiUrl = "http://metar.vatsim.net/metar.php?id=%s"
const timeout = time.Second * 10

func newApi() vatsimApi {
	return &api{
		log:    logging.New("modules/vatsim/api"),
		client: &http.Client{Timeout: timeout},
	}

}

type api struct {
	log    logging.Logger
	client *http.Client
}

func (a *api) GetMetar(icao string) (string, error) {
	url := fmt.Sprintf(metarApiUrl, icao)
	resp, err := a.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		a.log.Error("received a status code indicating an error", "url", url, "code", resp.StatusCode)
		return "", errors.New("http request error")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	a.log.Debug("received a response from the metar api", "text", string(body), "code", resp.StatusCode)
	text := strings.TrimSpace(string(body))
	return text, nil
}
