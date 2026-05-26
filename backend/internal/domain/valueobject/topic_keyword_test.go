package valueobject

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTopicKeyword(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      TopicKeyword
		wantError bool
	}{
		{name: "lowercase", input: "AI", want: TopicKeyword("ai"), wantError: false},
		{name: "already lowercase", input: "machine learning", want: TopicKeyword("machine learning"), wantError: false},
		{name: "trim spaces", input: "  blockchain  ", want: TopicKeyword("blockchain"), wantError: false},
		{name: "mixed case", input: "Deep Learning", want: TopicKeyword("deep learning"), wantError: false},
		{name: "turkish chars", input: "Yapay Zeka", want: TopicKeyword("yapay zeka"), wantError: false},
		{name: "empty string errors", input: "", wantError: true},
		{name: "only spaces errors", input: "   ", wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTopicKeyword(tt.input)
			if tt.wantError {
				assert.Error(t, err)
				assert.ErrorIs(t, err, ErrEmptyKeyword)
				assert.Empty(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestTopicKeyword_String(t *testing.T) {
	kw, err := NewTopicKeyword("AI")
	require.NoError(t, err)
	assert.Equal(t, "ai", kw.String())
}
