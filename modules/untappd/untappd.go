package untappd

import (
	"fmt"
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/modules"
	untappd2 "github.com/mdlayher/untappd"
	"strings"
	"net/http"
)

func New(sender modules.Sender, conf config.Config) (modules.Module, error) {
	client := http.Client{}
	utAPI, err := untappd2.NewClient(conf.UntappdClientID, conf.UntappdClientSecret, &client)
	if err != nil {
		return nil, err
	}

	rv := &untappd{
		Base: modules.NewBase("untappd", sender, conf),
		client: utAPI,
	}
	rv.AddCommand("ut", rv.ut)
	return rv, nil
}

type untappd struct {
	modules.Base
	client *untappd2.Client
}

func (u *untappd) ut(arguments modules.CommandArguments) ([]string, error) {
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
		fmt.Sprintf("%s | %s | %s | %.2f%% | IBU: %d | Rating: %.3f | !https://untappd.com/beer/%d",
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
