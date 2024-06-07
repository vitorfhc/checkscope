package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vitorfhc/checkscope/pkg/checkscope"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	scopeFileFlag := flag.String("f", "scope.txt", "Scope file")
	outOfScopeFileFlag := flag.String("o", "", "Out of scope file")
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

	outs := []string{}
	if *outOfScopeFileFlag != "" {
		log.Info().Msgf("Reading out of scope from %s", *outOfScopeFileFlag)
		outs, err = readFileLines(*outOfScopeFileFlag)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to read out of scope file")
		}
	}

	log.Info().Msg("Reading input from stdin")
	var inputs []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			inputs = append(inputs, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal().Err(err).Msg("Failed to read input from stdin")
	}

	cs := checkscope.New(inputs, scopes, outs)

	matches, err := cs.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to run checkscope")
	}

	typeToMatch := checkscope.MATCH_TYPE_IN_SCOPE
	if *reverseFlag {
		typeToMatch = checkscope.MATCH_TYPE_OUT_SCOPE
	}

	for _, match := range matches {
		if match.MatchType == typeToMatch {
			fmt.Println(match.Input)
		}
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
