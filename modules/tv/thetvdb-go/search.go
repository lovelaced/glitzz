package thetvdb

import (
	"net/http"
)

func newSearchService(base *baseService) *SearchService {
	rv := &SearchService{
		Url:         "search",
		UrlSeries:   "series",
		baseService: base,
	}
	return rv
}

type SearchService struct {
	Url       string
	UrlSeries string

	*baseService
}

type SearchSeriesParams struct {
	Name     *string
	ImdbId   *string
	Zap2itId *string
}

type SearchSeriesResult struct {
	Data []SeriesResult `json:"data"`
}

func (s *SearchService) Series(params SearchSeriesParams) (*SearchSeriesResult, *http.Response, error) {
	req, err := newRequest(http.MethodGet, s.Url, s.UrlSeries, nil)
	if err != nil {
		return nil, nil, err
	}
	setQuery(req, map[string]*string{
		"name":     params.Name,
		"imdbId":   params.ImdbId,
		"zap2itId": params.Zap2itId,
	})
	var result SearchSeriesResult
	resp, err := s.do(req, &result)
	if err != nil {
		return nil, resp, err
	}
	return &result, resp, nil
}
