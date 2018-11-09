package stackexchange

import (
	"github.com/PuffyVulva/Stack-on-Go/stackongo"
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"

	"errors"
	"html"
	"math/rand"
	"strings"
)

var (
	// Put your default tags here
	defaultTags []string
	lastLink    string
)

const (
	errInvalidSite = "Invalid site: "
	errInvalidTags = "Tag(s) not found"
	defaultSite    = "workplace"
	sePrefix       = "se"
	seLast         = "last"
)

// se = seprefix
// seprefix: seq (question title)
// seq: question title

// New initializes the stackexchange module
func New(sender core.Sender, conf config.Config) (core.Module, error) {
	rv := &stackexchange{
		Base: core.NewBase("stackexchange", sender, conf),
	}
	rv.AddCommand("so", rv.seStackOverflow)
	rv.AddCommand(sePrefix, rv.seTitle)
	rv.AddCommand(sePrefix+seLast, rv.seLastLink)
	return rv, nil
}

type stackexchange struct {
	core.Base
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
	return defaultTags
}

func (s *stackexchange) getRandQuestion(arguments core.CommandArguments) (*stackongo.Question, error) {
	site := s.getSite(arguments)
	tags := s.getTags(arguments)

	session := stackongo.NewSession(site)
	params := make(stackongo.Params)
	params.Add("sort", "creation")
	params.AddVectorized("tagged", tags)

	questions, err := session.AllQuestions(params)

	if err != nil {
		if strings.Contains(err.Error(), "No site found for name") {
			return nil, errors.New(errInvalidSite + site)
		}
		return nil, err
	}
	if len(questions.Items) == 0 {
		return nil, errors.New(errInvalidTags)
	}

	index := rand.Intn(len(questions.Items))
	return &questions.Items[index], nil
}

func (s *stackexchange) seTitle(arguments core.CommandArguments) ([]string, error) {
	question, err := s.getRandQuestion(arguments)

	if err != nil {
		if err.Error() == errInvalidTags || strings.Contains(err.Error(), errInvalidSite) {
			return []string{err.Error()}, nil
		}
		return nil, err
	}
	lastLink = question.Link
	return []string{html.UnescapeString(question.Title)}, err
}

func (s *stackexchange) seLastLink(arguments core.CommandArguments) ([]string, error) {
	return []string{lastLink}, nil
}

func (s *stackexchange) seStackOverflow(arguments core.CommandArguments) ([]string, error) {
	args := arguments
	args.Arguments = append([]string{"stackoverflow"}, args.Arguments...)
	return s.seTitle(args)
}
