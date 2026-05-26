package gateway

import (
	"testing"
)

func TestExtractField(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		prefix   string
		expected string
	}{
		{
			name:     "Standard field",
			text:     "BASLIK: Daha yavas ama daha iyi kod\nMETIN: Deneme",
			prefix:   "BASLIK:",
			expected: "Daha yavas ama daha iyi kod",
		},
		{
			name:     "Case insensitivity",
			text:     "baslik: Daha yavas ama daha iyi kod",
			prefix:   "BASLIK:",
			expected: "Daha yavas ama daha iyi kod",
		},
		{
			name:     "Turkish characters in text capitalized",
			text:     "BAŞLIK: Daha yavaş ama daha iyi kod",
			prefix:   "BASLIK:",
			expected: "Daha yavaş ama daha iyi kod",
		},
		{
			name:     "Missing prefix",
			text:     "METIN: Deneme",
			prefix:   "BASLIK:",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractField(tt.text, tt.prefix)
			if got != tt.expected {
				t.Errorf("extractField(%q, %q) = %q; want %q", tt.text, tt.prefix, got, tt.expected)
			}
		})
	}
}

func TestExtractFieldMultiLine(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		prefix   string
		expected string
	}{
		{
			name:     "Multi-line extract",
			text:     "OZET:\nBu bir test ozetidir.\nIkinci satir buradadir.",
			prefix:   "OZET:",
			expected: "Bu bir test ozetidir.\nIkinci satir buradadir.",
		},
		{
			name:     "Case insensitivity and trailing space",
			text:     "ozet:  Bu bir test ozetidir.  ",
			prefix:   "OZET:",
			expected: "Bu bir test ozetidir.",
		},
		{
			name:     "Missing prefix",
			text:     "Icerik buradadir.",
			prefix:   "OZET:",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractFieldMultiLine(tt.text, tt.prefix)
			if got != tt.expected {
				t.Errorf("extractFieldMultiLine(%q, %q) = %q; want %q", tt.text, tt.prefix, got, tt.expected)
			}
		})
	}
}
