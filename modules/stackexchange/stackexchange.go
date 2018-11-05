package stackexchange

import (
	"github.com/PuffyVulva/Stack-on-Go/stackongo"

	"errors"
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"math/rand"
	// "github.com/lovelaced/glitzz/logging"
)

var (
	// Put your default tags here
	defaultTags []string
	lastLink    string
)

const (
	errNoQuestion = "No Questions for given tags/site"
	defaultSite   = "stackoverflow"
	sePrefix      = "se"
	seAnswer      = "a"
	seLast        = "last"
)

// se = seprefix
// seprefix: seq (question title)
// seq: question title

// New initializes the stackexchange module
func New(sender core.Sender, conf config.Config) (core.Module, error) {
	rv := &stackexchange{
		Base: core.NewBase("stackexchange", sender, conf),
	}
	rv.AddCommand(sePrefix, rv.seTitle)
	rv.AddCommand(sePrefix+seAnswer, rv.seTitle)
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
		return nil, err
	}
	if len(questions.Items) == 0 {
		return nil, errors.New(errNoQuestion)
	}

	index := rand.Intn(len(questions.Items))
	return &questions.Items[index], nil
}

func (s *stackexchange) seTitle(arguments core.CommandArguments) ([]string, error) {
	question, err := s.getRandQuestion(arguments)

	if err != nil {
		if err.Error() == errNoQuestion {
			return []string{err.Error()}, nil
		}
		return nil, err
	}
	lastLink = question.Link
	return []string{question.Title}, err
}

func (s *stackexchange) seLastLink(arguments core.CommandArguments) ([]string, error) {
	return []string{lastLink}, nil
}
