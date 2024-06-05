package utils

import (
	"regexp"
	"strings"
)

// EscapeSpecialRegexChars escapes special regex characters.
// List of special characters: . + ? ^ $ ( ) [ ] { } |
func EscapeSpecialRegexChars(input string) string {
	specialChars := []string{".", "+", "?", "^", "$", "(", ")", "[", "]", "{", "}", "|"}
	for _, char := range specialChars {
		input = strings.ReplaceAll(input, char, "\\"+char)
	}
	return input
}

// WildcardToRegex converts a wildcard string to a regex string.
func WildcardToRegex(input string) string {
	r := regexp.MustCompile(`\*+`)
	input = r.ReplaceAllString(input, "*")

	// input = strings.ReplaceAll(input, "*.", "*")

	input = EscapeSpecialRegexChars(input)

	input = strings.ReplaceAll(input, "*", ".*")
	input = "^" + input + "$"

	return input
}
