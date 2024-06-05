package checkscope

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/vitorfhc/checkscope/pkg/utils"
)

type CheckScope struct {
	inputs   []string
	matchers []*regexp.Regexp
}

type Match struct {
	Input   string
	Scope   string
	Matched bool
}

func New(inputs []string, scopes []string) *CheckScope {
	cs := &CheckScope{
		inputs:   inputs,
		matchers: make([]*regexp.Regexp, 0),
	}

	for _, scope := range scopes {
		regexString := utils.WildcardToRegex(scope)
		compiled := regexp.MustCompile(regexString)
		cs.matchers = append(cs.matchers, compiled)
	}

	return cs
}

func (c *CheckScope) Run() ([]Match, error) {
	var matches []Match

	for _, input := range c.inputs {
		urlInput := input

		if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
			urlInput = "http://" + input
		}

		parsedUrl, err := url.Parse(urlInput)
		if err != nil {
			return nil, fmt.Errorf("failed to parse URL %q: %w", input, err)
		}

		match := Match{
			Input:   input,
			Scope:   "",
			Matched: false,
		}

		hostname := parsedUrl.Hostname()
		for _, matcher := range c.matchers {
			if matcher.MatchString(hostname) {
				match.Scope = matcher.String()
				match.Matched = true
				break
			}
		}

		matches = append(matches, match)
		log.Debug().Interface("match", match).Msg("Match result")
	}

	return matches, nil
}
