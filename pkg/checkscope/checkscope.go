package checkscope

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/vitorfhc/checkscope/pkg/utils"
)

const (
	MATCH_TYPE_IN_SCOPE  = "in_scope"
	MATCH_TYPE_OUT_SCOPE = "out_scope"
)

type CheckScope struct {
	inputs   []string
	matchers []*regexp.Regexp
	outs     []*regexp.Regexp
}

type Match struct {
	Input     string
	Matcher   string
	Matched   bool
	MatchType string
}

func New(inputs []string, scopes []string, outs []string) *CheckScope {
	cs := &CheckScope{
		inputs:   inputs,
		matchers: make([]*regexp.Regexp, 0),
		outs:     make([]*regexp.Regexp, 0),
	}

	for _, scope := range scopes {
		regexString := utils.WildcardToRegex(scope)
		compiled := regexp.MustCompile(regexString)
		cs.matchers = append(cs.matchers, compiled)
	}

	for _, out := range outs {
		regexString := utils.WildcardToRegex(out)
		compiled := regexp.MustCompile(regexString)
		cs.outs = append(cs.outs, compiled)
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
			Input:     input,
			Matcher:   "",
			Matched:   false,
			MatchType: MATCH_TYPE_OUT_SCOPE,
		}

		hostname := parsedUrl.Hostname()
		for _, matcher := range c.outs {
			if matcher.MatchString(hostname) {
				match.Matcher = matcher.String()
				match.Matched = true
				match.MatchType = MATCH_TYPE_OUT_SCOPE
				break
			}
		}

		if !match.Matched {
			for _, matcher := range c.matchers {
				if matcher.MatchString(hostname) {
					match.Matcher = matcher.String()
					match.Matched = true
					match.MatchType = MATCH_TYPE_IN_SCOPE
					break
				}
			}
		}

		matches = append(matches, match)
		log.Debug().Interface("match", match).Msg("Match result")
	}

	return matches, nil
}
