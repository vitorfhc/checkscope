package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	scopeFileFlag := flag.String("f", "scope.txt", "Scope file")
	debugFlag := flag.Bool("d", false, "Debug mode")
	silentFlag := flag.Bool("s", false, "Silent mode")
	reverseFlag := flag.Bool("r", false, "Reverse mode (prints out of scope URLs)")
	flag.Parse()

	if *scopeFileFlag == "" {
		log.Fatal().Msg("Scope file is required")
	}

	if *debugFlag {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if *silentFlag {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}

	log.Info().Msgf("Reading scope from %s", *scopeFileFlag)
	scopes, err := readFileLines(*scopeFileFlag)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read scope file")
	}

	var scopeCompiledRegex []*regexp.Regexp
	r, err := regexp.Compile(`\*+\.?`)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to compile regex")
	}
	for i, scope := range scopes {
		scopes[i] = r.ReplaceAllString(scope, "*")
		scopes[i] = wildcardToRegex(scopes[i])
		scopeCompiledRegex = append(scopeCompiledRegex, regexp.MustCompile(scopes[i]))
		log.Debug().Str("scope", scopes[i]).Msg("Scope regex added")
	}

	log.Info().Msg("Reading input from stdin")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		match := false
		line := scanner.Text()
		parsedUrl, err := url.Parse(line)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to parse URL")
		}
		hostname := parsedUrl.Hostname()
		for _, scope := range scopeCompiledRegex {
			if scope.MatchString(hostname) {
				match = true
				break
			}
		}

		if match && !*reverseFlag {
			fmt.Println(line)
		} else if !match && *reverseFlag {
			fmt.Println(line)
		}

		log.Debug().Str("line", line).Bool("match", match).Msg("Line processed")
	}
}

func readFileLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func escapeSpecialRegexChars(input string) string {
	specialChars := []string{".", "+", "?", "^", "$", "(", ")", "[", "]", "{", "}", "|"}
	for _, char := range specialChars {
		input = strings.ReplaceAll(input, char, "\\"+char)
	}
	return input
}

func wildcardToRegex(pattern string) string {
	// Escape special regex characters first
	escapedPattern := escapeSpecialRegexChars(pattern)
	// Replace wildcard * with regex .*
	regexPattern := strings.ReplaceAll(escapedPattern, "*", ".*")
	// Add anchors to match the whole string
	return "^" + regexPattern + "$"
}
