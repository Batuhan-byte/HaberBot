// Package valueobject defines domain value objects for the HaberBot application.
// Value objects are immutable types that represent descriptive aspects of the domain
// with no conceptual identity.
package valueobject

import "fmt"

// SourceType represents the origin source of an article.
type SourceType string

const (
	// SourceHackerNews represents articles fetched from the HackerNews API.
	SourceHackerNews SourceType = "hackernews"
	// SourceRSS represents articles fetched from RSS feeds.
	SourceRSS SourceType = "rss"
)

// validSourceTypes holds all valid source type values for validation.
var validSourceTypes = map[SourceType]bool{
	SourceHackerNews: true,
	SourceRSS:        true,
}

// IsValid checks whether the SourceType is a recognized, valid value.
func (s SourceType) IsValid() bool {
	return validSourceTypes[s]
}

// String returns the string representation of the SourceType.
func (s SourceType) String() string {
	return string(s)
}

// ParseSourceType converts a raw string into a validated SourceType.
// Returns an error if the string does not match a known source type.
func ParseSourceType(raw string) (SourceType, error) {
	source := SourceType(raw)
	if !source.IsValid() {
		return "", fmt.Errorf("invalid source type: %q", raw)
	}
	return source, nil
}
