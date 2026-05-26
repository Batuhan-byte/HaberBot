package fetcher

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	nethtml "golang.org/x/net/html"
)

// HTMLSanitizer is a deep module that parses raw HTML content,
// strips out dangerous tags (scripts, styles, iframes, form elements)
// and cleans custom attributes, ensuring the markup is safe and style-neutral.
type HTMLSanitizer struct{}

// NewHTMLSanitizer creates a new HTMLSanitizer.
func NewHTMLSanitizer() *HTMLSanitizer {
	return &HTMLSanitizer{}
}

// Sanitize parses a raw HTML string and returns a clean, safe, and style-neutral HTML markup.
func (s *HTMLSanitizer) Sanitize(html string) (string, error) {
	if html == "" {
		return "", nil
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}

	// 1. Remove dangerous/noisy tags
	doc.Find("script, style, iframe, link, meta, object, embed, input, form, button, textarea").Each(func(i int, sel *goquery.Selection) {
		sel.Remove()
	})

	// 2. Clean all attributes on remaining tags to enforce stylesheet neutrality
	doc.Find("*").Each(func(i int, sel *goquery.Selection) {
		node := sel.Get(0)
		var safeAttrs []nethtml.Attribute
		
		for _, attr := range node.Attr {
			key := strings.ToLower(attr.Key)
			
			// Retain only safe semantic attributes
			if key == "href" || key == "src" || key == "title" || key == "alt" || key == "rel" {
				safeAttrs = append(safeAttrs, attr)
			}
		}
		
		// Re-assign sanitized attributes
		node.Attr = safeAttrs
	})

	// 3. Extract cleaned body HTML
	var bodyHTML string
	bodySelection := doc.Find("body")
	if bodySelection.Length() > 0 {
		bodyHTML, _ = bodySelection.Html()
	} else {
		bodyHTML, _ = doc.Html()
	}

	return strings.TrimSpace(bodyHTML), nil
}
