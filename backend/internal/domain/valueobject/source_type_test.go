package valueobject

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSourceType_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		source   SourceType
		expected bool
	}{
		{name: "hackernews is valid", source: SourceHackerNews, expected: true},
		{name: "rss is valid", source: SourceRSS, expected: true},
		{name: "empty is invalid", source: SourceType(""), expected: false},
		{name: "unknown is invalid", source: SourceType("twitter"), expected: false},
		{name: "case-sensitive invalid", source: SourceType("HackerNews"), expected: false},
		{name: "spaces invalid", source: SourceType(" hackernews"), expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.source.IsValid())
		})
	}
}

func TestSourceType_String(t *testing.T) {
	assert.Equal(t, "hackernews", SourceHackerNews.String())
	assert.Equal(t, "rss", SourceRSS.String())
}

func TestParseSourceType(t *testing.T) {
	tests := []struct {
		name      string
		raw       string
		want      SourceType
		wantError bool
	}{
		{name: "parse hackernews", raw: "hackernews", want: SourceHackerNews, wantError: false},
		{name: "parse rss", raw: "rss", want: SourceRSS, wantError: false},
		{name: "empty string errors", raw: "", wantError: true},
		{name: "invalid source errors", raw: "invalid", wantError: true},
		{name: "case sensitive errors", raw: "HACKERNEWS", wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSourceType(tt.raw)
			if tt.wantError {
				assert.Error(t, err)
				assert.Empty(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
