package news

import (
	"regexp"
	"strings"
)

var htmlTagRegex *regexp.Regexp

func init() {
	htmlTagRegex = regexp.MustCompile(`<[^>]*>`)
}

func stripHtmlTags(text string) string {
	cleaned := htmlTagRegex.ReplaceAllString(text, "")

	cleaned = strings.Join(strings.Fields(cleaned), " ")
	return strings.TrimSpace(cleaned)
}
