package stackexchange

import (
	"github.com/ShikiCanKillServants/Stack-on-Go/stackongo"
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"

	"errors"
	"html"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const (
	errInvalidSite = "Invalid site: "
	errInvalidTags = "Tag(s) not found"
	defaultSite    = "workplace"
)

// New registers the stackexchange module
func New(sender core.Sender, conf config.Config) (core.Module, error) {
	rv := &stackexchange{
		Base: core.NewBase("stackexchange", sender, conf),
	}

	rv.AddCommand("so", rv.seStackOverflow)
	rv.AddCommand("se", rv.seTitle)
	rv.AddCommand("solast", rv.seLastSOLink)
	rv.AddCommand("selast", rv.seLastSELink)
	return rv, nil
}

type stackexchange struct {
	core.Base
	lastSELink string
	lastSOLink string
}

func (s *stackexchange) getSite(arguments core.CommandArguments) string {
	if len(arguments.Arguments) > 0 {
		return arguments.Arguments[0]
	}
	return defaultSite
}

func (s *stackexchange) getTags(arguments core.CommandArguments) []string {
	if len(arguments.Arguments) > 1 {
		return arguments.Arguments[1:]
	}
	return nil
}

func (s *stackexchange) getRandQuestion(arguments core.CommandArguments) (*stackongo.Question, string, error) {
	site := s.getSite(arguments)
	tags := s.getTags(arguments)

	netClient := http.Client{
		Timeout: time.Second * 10,
	}
	stackongo.SetHTTPClient(netClient)
	session := stackongo.NewSession(site)
	params := make(stackongo.Params)
	params.Add("sort", "creation")
	params.AddVectorized("tagged", tags)

	questions, err := session.AllQuestions(params)
	if err != nil {
		if strings.Contains(err.Error(), "No site found for name") {
			return nil, site, errors.New(errInvalidSite + site)
		}
		return nil, site, err
	}
	if len(questions.Items) == 0 {
		return nil, site, errors.New(errInvalidTags)
	}

	index := rand.Intn(len(questions.Items))
	return &questions.Items[index], site, nil
}

func (s *stackexchange) seTitle(arguments core.CommandArguments) ([]string, error) {
	question, site, err := s.getRandQuestion(arguments)
	if err != nil {
		if err.Error() == errInvalidTags || strings.Contains(err.Error(), errInvalidSite) {
			return []string{err.Error()}, nil
		}
		return nil, err
	}

	if site == "stackoverflow" {
		s.lastSOLink = question.Link
	} else {
		s.lastSELink = question.Link
	}
	return []string{html.UnescapeString(question.Title)}, err
}

func (s *stackexchange) seLastSOLink(arguments core.CommandArguments) ([]string, error) {
	if s.lastSOLink == "" {
		return []string{"https://stackoverflow.com"}, nil
	}
	return []string{s.lastSOLink}, nil
}

func (s *stackexchange) seLastSELink(arguments core.CommandArguments) ([]string, error) {
	if s.lastSELink == "" {
		return []string{"https://stackexchange.com"}, nil
	}
	return []string{s.lastSELink}, nil
}

func (s *stackexchange) seStackOverflow(arguments core.CommandArguments) ([]string, error) {
	args := arguments
	args.Arguments = append([]string{"stackoverflow"}, args.Arguments...)
	return s.seTitle(args)
}
