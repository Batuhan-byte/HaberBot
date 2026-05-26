package valueobject

import (
	"errors"
	"strings"
)

// TopicKeyword represents a single keyword used to filter and match articles
// for a given topic. Keywords are always stored in lowercase trimmed form.
type TopicKeyword string

// ErrEmptyKeyword is returned when attempting to create a keyword from an empty string.
var ErrEmptyKeyword = errors.New("topic keyword cannot be empty")

// NewTopicKeyword creates a new TopicKeyword from the given value.
// It trims whitespace and validates that the result is non-empty.
func NewTopicKeyword(value string) (TopicKeyword, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", ErrEmptyKeyword
	}
	return TopicKeyword(strings.ToLower(trimmed)), nil
}

// String returns the string representation of the TopicKeyword.
func (k TopicKeyword) String() string {
	return string(k)
}
