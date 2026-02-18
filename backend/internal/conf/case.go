package conf

import (
	"strings"
	"unicode"
)

func toConstantCase(fieldName string) string {
	return toCase(fieldName, strings.ToUpper, "_")
}

func toKebabCase(fieldName string) string {
	return toCase(fieldName, strings.ToLower, "-")
}

func toCase(fieldName string, caseFunc func(str string) string, delim string) string {
	split := splitCamelCase(fieldName)
	for i, word := range split {
		split[i] = caseFunc(word)
	}
	return strings.Join(split, delim)
}

func splitCamelCase(s string) []string {
	runes := []rune(s)
	var words []string
	lastStart := 0

	for i := 1; i < len(runes); i++ {
		prev := runes[i-1]
		curr := runes[i]

		split := false

		// lower -> upper
		if unicode.IsLower(prev) && unicode.IsUpper(curr) {
			split = true
		}

		// letter <-> digit
		if unicode.IsLetter(prev) && unicode.IsDigit(curr) ||
			unicode.IsDigit(prev) && unicode.IsLetter(curr) {
			split = true
		}

		// upper -> upper -> lower (acronym end)
		if i >= 2 &&
			unicode.IsUpper(runes[i-2]) &&
			unicode.IsUpper(prev) &&
			unicode.IsLower(curr) {
			words = append(words, string(runes[lastStart:i-1]))
			lastStart = i - 1
			continue
		}

		if split {
			words = append(words, string(runes[lastStart:i]))
			lastStart = i
		}
	}

	if lastStart < len(runes) {
		words = append(words, string(runes[lastStart:]))
	}

	return words
}
