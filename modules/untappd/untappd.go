package untappd

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	untappd2 "github.com/mdlayher/untappd"
	"github.com/pkg/errors"
)

func New(sender core.Sender, conf config.Config) (core.Module, error) {
	base := core.NewBase("untappd", sender, conf)
	if configIsIncorrect(conf) {
		return nil, errors.New("Invalid config, see: https://untappd.com/api/docs")
	}
	client := http.Client{Timeout: 20 * time.Second}
	utAPI, err := untappd2.NewClient(conf.Untappd.ClientID, conf.Untappd.ClientSecret, &client)
	if err != nil {
		return nil, errors.Wrap(err, "creating untappd client failed")
	}
	rv := &untappd{
		Base:   base,
		client: utAPI,
	}
	rv.AddCommand("ut", rv.ut)
	return rv, nil
}

type untappd struct {
	core.Base
	client *untappd2.Client
}

func (u *untappd) ut(arguments core.CommandArguments) ([]string, error) {
	beerResults, _, err := u.client.Beer.Search(strings.Join(arguments.Arguments, " "))
	if len(beerResults) < 1 {
		return []string{"No beers found"}, err
	}
	rawInfo, _, err := u.client.Beer.Info(beerResults[0].ID, true)
	if err != nil {
		return nil, err
	}
	u.Log.Debug("untappdAPI", "beer", beerResults[0])
	var text = []string{
		fmt.Sprintf("%s | %s | %s | %.2f%% ABV | IBU: %d | Rating: %.3f | !https://untappd.com/beer/%d",
			rawInfo.Name,
			rawInfo.Style,
			rawInfo.Brewery.Name,
			rawInfo.ABV,
			rawInfo.IBU,
			rawInfo.OverallRating,
			rawInfo.ID),
	}
	return text, err
}

func configIsIncorrect(conf config.Config) bool {
	return conf.Untappd.ClientID == "" ||
		conf.Untappd.ClientSecret == "" ||
		conf.Untappd.ClientID == config.Default().Untappd.ClientID ||
		conf.Untappd.ClientSecret == config.Default().Untappd.ClientSecret
}
