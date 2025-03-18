package post

import "regexp"

var (
	contentRegex = regexp.MustCompile(`(?s)^\s*\S.{0,299}$`)
)

func validateContent(content string) error {
	if !contentRegex.MatchString(content) {
		return ErrInvalidContent
	}
	return nil
}
