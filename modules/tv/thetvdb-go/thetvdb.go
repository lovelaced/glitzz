package thetvdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const ApiVersion = "2.2.0"

type requestPerformer interface {
	do(req *http.Request) (*http.Response, error)
	doUnauthorized(req *http.Request) (*http.Response, error)
}

func NewClient(loginParams LoginParams, httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	client := &Client{
		Url: "https://api.thetvdb.com",

		httpClient:  httpClient,
		loginParams: loginParams,
	}
	base := &baseService{r: client}
	client.Search = newSearchService(base)
	client.Series = newSeriesService(base)
	client.login = newLoginService(base)
	return client, nil
}

type Client struct {
	Search *SearchService
	Series *SeriesService
	login  *loginService

	Url string

	token       *string
	httpClient  *http.Client
	loginParams LoginParams
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	if c.token == nil {
		err := c.authenticate(c.loginParams)
		if err != nil {
			return nil, errors.Wrap(err, "authentication failed")
		}
	}
	resp, err := c.internalDo(req)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	if resp.StatusCode == http.StatusUnauthorized {
		err := c.authenticate(c.loginParams)
		if err != nil {
			return nil, errors.Wrap(err, "authentication failed")
		}
		return c.internalDo(req)
	}
	return resp, nil
}

func (c *Client) doUnauthorized(req *http.Request) (*http.Response, error) {
	return c.internalDo(req)
}

func (c *Client) internalDo(req *http.Request) (*http.Response, error) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", fmt.Sprintf("application/vnd.thetvdb.v%s", ApiVersion))
	if c.token != nil {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *c.token))
	}
	u, err := joinUrl(c.Url, req.URL.String())
	if err != nil {
		return nil, errors.Wrap(err, "could not join urls")
	}
	req.URL, err = url.Parse(u)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse the url")
	}
	println(req.URL.String())
	resp, err := c.httpClient.Do(req)
	return resp, err
}

func (c *Client) authenticate(params LoginParams) error {
	result, _, err := c.login.Login(params)
	if err != nil {
		return errors.Wrap(err, "could not authenticate")
	}
	c.token = &result.Token
	return nil
}

type baseService struct {
	r requestPerformer
}

func (b *baseService) do(req *http.Request, target interface{}) (*http.Response, error) {
	resp, err := b.r.do(req)
	if err != nil {
		return nil, err
	}
	return b.unmarshal(resp, target)
}

func (b *baseService) doUnauthorized(req *http.Request, target interface{}) (*http.Response, error) {
	resp, err := b.r.doUnauthorized(req)
	if err != nil {
		return nil, err
	}
	return b.unmarshal(resp, target)
}

func (b *baseService) unmarshal(resp *http.Response, target interface{}) (*http.Response, error) {
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("received status %d", resp.StatusCode)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}
	return resp, json.Unmarshal(bodyBytes, &target)
}

func joinUrl(a, b string) (string, error) {
	var sep string
	if a != "" && b != "" {
		sep = "/"
	}
	return a + sep + b, nil
}

func newRequest(method, url1, url2 string, body interface{}) (*http.Request, error) {
	var buf io.Reader
	if body != nil {
		j, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(j)
	}
	url, err := joinUrl(url1, url2)
	if err != nil {
		return nil, err
	}
	return http.NewRequest(method, url, buf)
}

func setQuery(req *http.Request, query map[string]*string) {
	q := req.URL.Query()
	for key, value := range query {
		if value != nil {
			q.Add(key, *value)
		}
	}
	req.URL.RawQuery = q.Encode()
}

func intPtrToStrPtr(i *int) *string {
	if i == nil {
		return nil
	}
	s := strconv.Itoa(*i)
	return &s
}
