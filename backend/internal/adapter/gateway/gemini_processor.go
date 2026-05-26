// Package gateway provides external service adapters (gateways) that
// implement domain ports for third-party API integrations.
package gateway

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const geminiModelName = "gemini-2.5-flash"
const geminiTemperature = 0.3
const geminiMaxTokens = 2048
const geminiMaxContentLength = 30000 // Allow large articles for translation

// GeminiProcessor implements port.AIProcessor using Google's Generative AI SDK.
type GeminiProcessor struct {
	apiKey string
}

// NewGeminiProcessor creates a new GeminiProcessor.
func NewGeminiProcessor(apiKey string) *GeminiProcessor {
	return &GeminiProcessor{apiKey: apiKey}
}

func (p *GeminiProcessor) Translate(ctx context.Context, title, content string) (string, string, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(p.apiKey))
	if err != nil {
		return "", "", fmt.Errorf("creating gemini client: %w", err)
	}
	defer client.Close()

	prompt := fmt.Sprintf(translatePromptTemplate, title, truncateContent(content))
	response, err := p.generateContent(ctx, client, prompt)
	if err != nil {
		return "", "", err
	}

	turkishTitle := extractField(response, "BASLIK:")
	if turkishTitle == "" {
		turkishTitle = extractField(response, "TITLE:") 
	}

	turkishContent := extractFieldMultiLine(response, "METIN:")
	if turkishContent == "" {
		turkishContent = response // Fallback: just return everything if format fails
	}

	return turkishTitle, turkishContent, nil
}

func (p *GeminiProcessor) Summarize(ctx context.Context, content string) (string, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(p.apiKey))
	if err != nil {
		return "", fmt.Errorf("creating gemini client: %w", err)
	}
	defer client.Close()

	prompt := fmt.Sprintf(summarizePromptTemplate, truncateContent(content))
	response, err := p.generateContent(ctx, client, prompt)
	if err != nil {
		return "", err
	}

	summary := extractFieldMultiLine(response, "OZET:")
	if summary == "" {
		summary = response
	}

	return summary, nil
}

func (p *GeminiProcessor) generateContent(ctx context.Context, client *genai.Client, prompt string) (string, error) {
	model := client.GenerativeModel(geminiModelName)
	temperature := float32(geminiTemperature)
	model.Temperature = &temperature
	maxTokens := int32(geminiMaxTokens)
	model.MaxOutputTokens = &maxTokens

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("generating content: %w", err)
	}

	return extractTextFromResponse(resp)
}

func extractTextFromResponse(resp *genai.GenerateContentResponse) (string, error) {
	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates in response")
	}
	candidate := resp.Candidates[0]
	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return "", fmt.Errorf("no content parts in candidate")
	}
	var text string
	for _, part := range candidate.Content.Parts {
		if textPart, ok := part.(genai.Text); ok {
			text += string(textPart)
		}
	}
	return text, nil
}
