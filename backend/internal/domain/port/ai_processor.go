package port

import "context"

// AIProcessor defines the contract for artificial intelligence operations.
// Adapters implementing this interface may use OpenAI, Google Gemini, Anthropic, etc.
type AIProcessor interface {
	// Translate sends article content to the AI model for Turkish translation.
	// It returns the translated title and the full translated content.
	Translate(ctx context.Context, title, content string) (string, string, error)

	// Summarize generates a short 3-4 sentence summary from the given text.
	Summarize(ctx context.Context, content string) (string, error)

	// TranslateAndSummarize translates both title and content to Turkish and generates a short 3-4 sentence summary in a single call.
	TranslateAndSummarize(ctx context.Context, title, content string) (string, string, string, error)
}
