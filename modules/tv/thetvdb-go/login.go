package thetvdb

import (
	"net/http"
)

func newLoginService(base *baseService) *loginService {
	rv := &loginService{
		Url:             "",
		UrlLogin:        "login",
		UrlRefreshToken: "refresh_token",
		baseService:     base,
	}
	return rv
}

type loginService struct {
	Url             string
	UrlLogin        string
	UrlRefreshToken string

	*baseService
}

type LoginParams struct {
	ApiKey   string  `json:"apikey"`
	UserKey  *string `json:"userkey,omitempty"`
	UserName *string `json:"username,omitempty"`
}

type loginResult struct {
	Token string `json:"token"`
}

func (l *loginService) Login(params LoginParams) (*loginResult, *http.Response, error) {
	req, err := newRequest(http.MethodPost, l.Url, l.UrlLogin, params)
	if err != nil {
		return nil, nil, err
	}
	var result loginResult
	resp, err := l.doUnauthorized(req, &result)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp, nil
}
