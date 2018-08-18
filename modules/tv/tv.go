package tv

import (
	"fmt"
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"github.com/lovelaced/glitzz/modules/tv/thetvdb-go"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"time"
)

const timeout = 20 * time.Second

func New(sender core.Sender, conf config.Config) (core.Module, error) {
	if configIsIncorrect(conf) {
		return nil, errors.New("Invalid config, login on TheTVDB and see: https://www.thetvdb.com/member/api")
	}
	loginParams := thetvdb.LoginParams{
		ApiKey: conf.Tv.ApiKey,
	}
	httpClient := &http.Client{Timeout: timeout}
	api, err := thetvdb.NewClient(loginParams, httpClient)
	if err != nil {
		return nil, errors.Wrap(err, "could not create thetvdb api client")
	}
	rv := &tv{
		Base: core.NewBase("tv", sender, conf),
		api:  api,
	}
	rv.AddCommand("next_episode", rv.next_episode)
	return rv, nil
}

type tv struct {
	core.Base
	api *thetvdb.Client
}

func (t *tv) next_episode(arguments core.CommandArguments) ([]string, error) {
	if len(arguments.Arguments) > 0 {
		seriesName := strings.Join(arguments.Arguments, " ")
		text, err := t.getNextEpisodeText(seriesName)
		if err != nil {
			if err == seriesNotFoundError {
				return []string{"Series not found."}, nil
			}
			if err == nextEpisodeNotFoundError {
				return []string{"Next episode not found."}, nil
			}
			return nil, errors.Wrap(err, "could not find the next episode")
		}
		return []string{text}, nil
	}
	return nil, nil
}

func (t *tv) getNextEpisodeText(seriesName string) (string, error) {
	seriesId, err := t.getSeriesId(seriesName)
	if err != nil {
		return "", err
	}
	episodes, err := t.getAllEpisodes(seriesId)
	if err != nil {
		return "", err
	}
	episode := t.findNextEpisode(episodes)
	if episode == nil {
		return "", nextEpisodeNotFoundError
	}
	text := fmt.Sprintf("Episode %dx%d airs on %s.", *episode.AiredSeason, *episode.AiredEpisodeNumber, *episode.FirstAired)
	return text, nil
}

var seriesNotFoundError = errors.New("series not found")
var nextEpisodeNotFoundError = errors.New("next episode not found")

func (t *tv) findNextEpisode(episodes []thetvdb.EpisodeResult) *thetvdb.EpisodeResult {
	var closestEpisode *thetvdb.EpisodeResult
	for _, episode := range episodes {
		if episode.AiredSeason == nil || episode.AiredEpisodeNumber == nil ||
			episode.FirstAired == nil {
			continue
		}
		t.Log.Debug("processing episode", "season", *episode.AiredSeason, "number", *episode.AiredEpisodeNumber, "first_aired", *episode.FirstAired)
		episodeDate, err := parseDate(episode.FirstAired)
		if err != nil {
			t.Log.Debug("could not parse date", "err", err, "date", *episode.FirstAired)
			continue
		}
		if episodeDate.Before(time.Now()) {
			continue
		}
		if closestEpisode == nil {
			episodeCopy := episode
			closestEpisode = &episodeCopy
		} else {
			closestEpisodeDate, err := parseDate(closestEpisode.FirstAired)
			if err != nil {
				t.Log.Debug("could not parse date", "err", err, "date", *episode.FirstAired)
				continue
			}
			if closestEpisodeDate.After(episodeDate) {
				episodeCopy := episode
				closestEpisode = &episodeCopy
			}
		}
	}
	return closestEpisode
}

func (t *tv) getSeriesId(name string) (int, error) {
	params := thetvdb.SearchSeriesParams{
		Name: &name,
	}
	series, resp, err := t.api.Search.Series(params)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return 0, seriesNotFoundError
		}
		return 0, err
	}
	if len(series.Data) == 0 {
		return 0, seriesNotFoundError
	}
	seriesResult := series.Data[0]
	if seriesResult.Id == nil {
		return 0, errors.New("series id was null")
	}
	return *seriesResult.Id, nil
}

func (t *tv) getAllEpisodes(seriesId int) ([]thetvdb.EpisodeResult, error) {
	var episodes []thetvdb.EpisodeResult
	nextPage := 1
	for {
		params := thetvdb.SeriesEpisodesParams{
			Id:   seriesId,
			Page: &nextPage,
		}
		episodesResult, _, err := t.api.Series.Episodes(params)
		if err != nil {
			return nil, err
		}
		episodes = append(episodes, episodesResult.Data...)
		if episodesResult.Links == nil {
			break
		}
		if episodesResult.Links.Next == nil {
			break
		}
		nextPage = *episodesResult.Links.Next
	}
	return episodes, nil
}

const dateFormat = "2006-01-02"

func parseDate(s *string) (time.Time, error) {
	if s == nil {
		return time.Time{}, errors.New("date was null")
	}
	return time.Parse(dateFormat, *s)
}

func configIsIncorrect(conf config.Config) bool {
	return conf.Tv.ApiKey == "" ||
		conf.Tv.ApiKey == config.Default().Tv.ApiKey
}
