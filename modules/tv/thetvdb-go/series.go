package thetvdb

import (
	"net/http"
	"strconv"
	"strings"
)

func newSeriesService(base *baseService) *SeriesService {
	rv := &SeriesService{
		Url:         "series",
		UrlGet:      "<id>",
		UrlEpisodes: "<id>/episodes",
		baseService: base,
	}
	return rv
}

type SeriesService struct {
	Url         string
	UrlGet      string
	UrlEpisodes string

	*baseService
}

type SeriesGetParams struct {
	Id int
}

type SeriesGetResult struct {
	Data SeriesResult `json:"data"`
}

type SeriesResult struct {
	Added           *string  `json:"added"`
	AirsDayOfWeek   *string  `json:"airsDayOfWeek"`
	AirsTime        *string  `json:"airsTime"`
	Aliases         []string `json:"aliases"`
	Banner          *string  `json:"banner"`
	FirstAired      *string  `json:"firstAired"`
	Genra           []string `json:"genre"`
	Id              *int     `json:"id"`
	ImdbId          *string  `json:"imdbId"`
	LastUpdated     *int     `json:"lastUpdated"`
	Network         *string  `json:"network"`
	NetworkId       *string  `json:"networkId"`
	Overview        *string  `json:"overview"`
	Rating          *string  `json:"rating"`
	Runtime         *string  `json:"runtime"`
	SeriesId        *string  `json:"seriesId"`
	SeriesName      *string  `json:"seriesName"`
	SiteRating      *float64 `json:"siteRating"`
	SiteRatingCount *int     `json:"siteRatingCount"`
	Slug            *string  `json:"slug"`
	Status          *string  `json:"status"`
	Zap2itId        *string  `json:"zap2itId"`
}

func (s *SeriesService) Get(params SeriesGetParams) (*SeriesGetResult, *http.Response, error) {
	url := strings.Replace(s.UrlGet, "<id>", strconv.Itoa(params.Id), -1)
	req, err := newRequest(http.MethodGet, s.Url, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var result SeriesGetResult
	resp, err := s.do(req, &result)
	if err != nil {
		return nil, resp, err
	}
	return &result, resp, nil
}

type SeriesEpisodesParams struct {
	Id   int
	Page *int // starts at 1
}

type SeriesEpisodesResult struct {
	Data  []EpisodeResult `json:"data"`
	Links *Links          `json:"links"`
}

type Links struct {
	First    *int `json:"first"`
	Last     *int `json:"last"`
	Next     *int `json:"next"`
	Previous *int `json:"previous"`
}

type EpisodeResult struct {
	AbsoluteNumber     *int     `json:"absoluteNumber"`
	AiredEpisodeNumber *int     `json:"airedEpisodeNumber"`
	AiredSeason        *int     `json:"airedSeason"`
	AirsAfterSeason    *int     `json:"airsAfterSeason"`
	AirsBeforeEpisode  *int     `json:"airsBeforeEpisode"`
	AirsBeforeSeason   *int     `json:"airsBeforeSeason"`
	Directors          []string `json:"directors"`
	DvdChapter         *int     `json:"dvdChapter"`
	DvdDiscid          *string  `json:"dvdDiscid"`
	DvdEpisodeNumber   *int     `json:"dvdEpisodeNumber"`
	DvdSeason          *int     `json:"dvdSeason"`
	EpisodeName        *string  `json:"episodeName"`
	Filename           *string  `json:"filename"`
	FirstAired         *string  `json:"firstAired"`
	GuestStars         []string `json:"guestStars"`
	Id                 *int     `json:"id"`
	ImdbId             *string  `json:"imdbId"`
	LastUpdated        *int     `json:"lastUpdated"`
	LastUpdatedBy      *int     `json:"lastUpdatedBy"`
	Overview           *string  `json:"overview"`
	ProductionCode     *string  `json:"productionCode"`
	SeriesId           *int     `json:"seriesId"`
	ShowUrl            *string  `json:"showUrl"`
	SiteRating         *float64 `json:"siteRating"`
	SiteRatingCount    *int     `json:"siteRatingCount"`
	ThumbAdded         *string  `json:"thumbAdded"`
	ThumbAuthor        *int     `json:"thumbAuthor"`
	ThumbHeight        *string  `json:"thumbHeight"`
	ThumbWidth         *string  `json:"thumbWidth"`
	Writers            []string `json:"writers"`
}

func (s *SeriesService) Episodes(params SeriesEpisodesParams) (*SeriesEpisodesResult, *http.Response, error) {
	url := strings.Replace(s.UrlEpisodes, "<id>", strconv.Itoa(params.Id), -1)
	req, err := newRequest(http.MethodGet, s.Url, url, nil)
	if err != nil {
		return nil, nil, err
	}
	setQuery(req, map[string]*string{
		"page": intPtrToStrPtr(params.Page),
	})
	var result SeriesEpisodesResult
	resp, err := s.do(req, &result)
	if err != nil {
		return nil, resp, err
	}
	return &result, resp, nil
}
