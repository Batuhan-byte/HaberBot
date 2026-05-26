package fetcher

import (
	"testing"
)

func TestHTMLSanitizer_Sanitize(t *testing.T) {
	sanitizer := NewHTMLSanitizer()

	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "Boş HTML girişi",
			input:    "",
			expected: "",
			wantErr:  false,
		},
		{
			name:     "script etiketini temizleme",
			input:    `<div>Hello <script>alert("xss")</script>World</div>`,
			expected: `<div>Hello World</div>`,
			wantErr:  false,
		},
		{
			name:     "style etiketini temizleme",
			input:    `<div>Hello <style>body { background: black; }</style>World</div>`,
			expected: `<div>Hello World</div>`,
			wantErr:  false,
		},
		{
			name:     "iframe ve form temizleme",
			input:    `<div><iframe src="http://evil.com"></iframe><form><input type="text"/></form>World</div>`,
			expected: `<div>World</div>`,
			wantErr:  false,
		},
		{
			name:     "Güvensiz nitelikleri temizleme (style, class, onclick)",
			input:    `<p style="color: red;" class="text-large" onclick="doEvil()">Hello World</p>`,
			expected: `<p>Hello World</p>`,
			wantErr:  false,
		},
		{
			name:     "Güvenli nitelikleri koruma (href, src, title, alt, rel)",
			input:    `<a href="https://example.com" title="Example Link" target="_blank" rel="nofollow">Link</a>`,
			expected: `<a href="https://example.com" rel="nofollow" title="Example Link">Link</a>`,
			wantErr:  false,
		},
		{
			name:     "Karmaşık HTML yapısını temizleme",
			input:    `<html><head><title>Ignore Me</title><style>p {color: blue;}</style></head><body><article><h1 class="title">Main Title</h1><p onclick="alert('click')">Some <script>console.log('x')</script>text <a href="/path">here</a>.</p></article></body></html>`,
			expected: `<article><h1>Main Title</h1><p>Some text <a href="/path">here</a>.</p></article>`,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sanitizer.Sanitize(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("HTMLSanitizer.Sanitize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("HTMLSanitizer.Sanitize() = %q, expected %q", got, tt.expected)
			}
		})
	}
}
