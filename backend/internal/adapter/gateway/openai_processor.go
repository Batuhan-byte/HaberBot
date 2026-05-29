package gateway

import (
	"context"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

const defaultTemperature = 0.3
const maxContentLength = 30000 // Allow large articles for translation

const translatePromptTemplate = `Sen bir teknoloji haberleri çevirmenisin.
Görevin: Verilen İngilizce makale başlığını ve metnini profesyonel, akıcı bir Türkçeye çevirmek.
Kurallar:
- Teknik terimler İngilizce kalmalı (ör: API, machine learning, framework, React, Go)
- Emoji kullanma
- Doğal ve akıcı Türkçe yaz
- SADECE çeviriyi ver, kendi yorumunu katma.
- İçerikteki tüm HTML etiketlerini (ör: <p>, <a>, <strong>, <ul>, <li>, <h2> vb.) kesinlikle koru. Çeviri yaparken etiketlerin yerlerini ve yapısını bozma, sadece etiketlerin içindeki metinleri çevir.
- Haberin başında veya sonunda yer alan; reklamlar, çerez onay metinleri, abonelik/paywall uyarıları (örneğin: "Skip Ad", "You are viewing a preview", "Sign in to continue", "Subscribe to read", "Erişim onaylandığında makale yüklenecektir", "ReklamATLA" vb.) gibi haber içeriğiyle ilgisi olmayan sistem, reklam veya üyelik mesajlarını TESPİT EDİP ÇIKAR. Bunları çeviriye dahil etme, doğrudan atla ve sadece asıl haberi çevir.

Başlık: %s
İçerik: %s

Yanıtını tam olarak şu formatta ver:
BASLIK: [Türkçe başlık]
METIN: [Türkçe tam çeviri]`

const summarizePromptTemplate = `Sen bir özetleyicisin.
Görevin: Verilen Türkçe makaleden 3-4 cümlelik "hap bilgi" özeti oluştur.
SADECE özeti ver.

İçerik: %s

Yanıtını tam olarak şu formatta ver:
OZET: [Özet metni]`

const translateAndSummarizePromptTemplate = `Sen bir teknoloji haberleri çevirmenisin.
Görevin: Verilen İngilizce makale başlığını ve metnini profesyonel, akıcı bir Türkçeye çevirmek ve aynı zamanda 3-4 cümlelik Türkçe "hap bilgi" özeti oluşturmak.
Kurallar:
- Teknik terimler İngilizce kalmalı (ör: API, machine learning, framework, React, Go)
- Emoji kullanma
- Doğal ve akıcı Türkçe yaz
- SADECE çeviriyi ve özeti ver, kendi yorumunu katma.
- İçerikteki tüm HTML etiketlerini (ör: <p>, <a>, <strong>, <ul>, <li>, <h2> vb.) kesinlikle koru. Çeviri yaparken etiketlerin yerlerini ve yapısını bozma, sadece etiketlerin içindeki metinleri çevir.
- Haberin başında veya sonunda yer alan; reklamlar, çerez onay metinleri, abonelik/paywall uyarıları (örneğin: "Skip Ad", "You are viewing a preview", "Sign in to continue", "Subscribe to read", "Erişim onaylandığında makale yüklenecektir", "ReklamATLA" vb.) gibi haber içeriğiyle ilgisi olmayan sistem, reklam veya üyelik mesajlarını TESPİT EDİP ÇIKAR. Bunları çeviriye dahil etme, doğrudan atla ve sadece asıl haberi çevir.

Başlık: %s
İçerik: %s

Yanıtını tam olarak şu formatta ver:
BASLIK: [Türkçe başlık]
METIN: [Türkçe tam çeviri]
OZET: [Özet metni]`

// OpenAIProcessor implements port.AIProcessor using an OpenAI-compatible client API.
type OpenAIProcessor struct {
	apiKey    string
	baseURL   string
	modelName string
}

// NewOpenAIProcessor creates a new OpenAI-compatible AI processor gateway.
func NewOpenAIProcessor(apiKey, baseURL, modelName string) *OpenAIProcessor {
	return &OpenAIProcessor{
		apiKey:    apiKey,
		baseURL:   baseURL,
		modelName: modelName,
	}
}

func (p *OpenAIProcessor) Translate(ctx context.Context, title, content string) (string, string, error) {
	config := openai.DefaultConfig(p.apiKey)
	config.BaseURL = p.baseURL
	client := openai.NewClientWithConfig(config)

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

func (p *OpenAIProcessor) Summarize(ctx context.Context, content string) (string, error) {
	config := openai.DefaultConfig(p.apiKey)
	config.BaseURL = p.baseURL
	client := openai.NewClientWithConfig(config)

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

func (p *OpenAIProcessor) TranslateAndSummarize(ctx context.Context, title, content string) (string, string, string, error) {
	config := openai.DefaultConfig(p.apiKey)
	config.BaseURL = p.baseURL
	client := openai.NewClientWithConfig(config)

	prompt := fmt.Sprintf(translateAndSummarizePromptTemplate, title, truncateContent(content))
	response, err := p.generateContent(ctx, client, prompt)
	if err != nil {
		return "", "", "", err
	}

	turkishTitle := extractField(response, "BASLIK:")
	if turkishTitle == "" {
		turkishTitle = extractField(response, "TITLE:") 
	}

	turkishContent := extractFieldBetween(response, "METIN:", "OZET:")
	if turkishContent == "" {
		turkishContent = extractFieldMultiLine(response, "METIN:")
	}
	if turkishContent == "" {
		turkishContent = response // Fallback: just return everything if format fails
	}

	turkishSummary := extractFieldMultiLine(response, "OZET:")

	return turkishTitle, turkishContent, turkishSummary, nil
}


func (p *OpenAIProcessor) generateContent(ctx context.Context, client *openai.Client, prompt string) (string, error) {
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: p.modelName,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: defaultTemperature,
	})
	if err != nil {
		return "", fmt.Errorf("generating content: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned in response")
	}

	return resp.Choices[0].Message.Content, nil
}

func cleanPrefix(s string) string {
	s = strings.ToUpper(s)
	s = strings.ReplaceAll(s, "Ş", "S")
	s = strings.ReplaceAll(s, "Ç", "C")
	s = strings.ReplaceAll(s, "Ğ", "G")
	s = strings.ReplaceAll(s, "İ", "I")
	s = strings.ReplaceAll(s, "Ö", "O")
	s = strings.ReplaceAll(s, "Ü", "U")
	s = strings.ReplaceAll(s, ":", "")
	return strings.TrimSpace(s)
}

func extractField(text, prefix string) string {
	lines := strings.Split(text, "\n")
	prefixClean := cleanPrefix(prefix)
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		colonIdx := strings.Index(trimmed, ":")
		if colonIdx == -1 {
			continue
		}
		
		key := strings.TrimSpace(trimmed[:colonIdx])
		val := strings.TrimSpace(trimmed[colonIdx+1:])
		
		if cleanPrefix(key) == prefixClean {
			return val
		}
	}
	
	// Fallback to searching variations in the whole text if not found line-by-line
	var variations []string
	if prefixClean == "BASLIK" {
		variations = []string{"BASLIK:", "BAŞLIK:", "TITLE:", "Baslik:", "Başlık:", "Title:", "baslik:", "başlık:", "title:"}
	} else if prefixClean == "OZET" {
		variations = []string{"OZET:", "ÖZET:", "SUMMARY:", "Ozet:", "Özet:", "Summary:", "ozet:", "özet:", "summary:"}
	} else {
		variations = []string{prefix}
	}
	
	for _, v := range variations {
		idx := strings.Index(text, v)
		if idx != -1 {
			rem := text[idx+len(v):]
			endOfLine := strings.Index(rem, "\n")
			if endOfLine != -1 {
				return strings.TrimSpace(rem[:endOfLine])
			}
			return strings.TrimSpace(rem)
		}
	}
	
	return ""
}

// extractFieldMultiLine extracts everything after the given prefix across multiple lines.
func extractFieldMultiLine(text, prefix string) string {
	prefixClean := cleanPrefix(prefix)
	
	var variations []string
	if prefixClean == "OZET" {
		variations = []string{"OZET:", "ÖZET:", "SUMMARY:", "Ozet:", "Özet:", "Summary:", "ozet:", "özet:", "summary:"}
	} else if prefixClean == "METIN" {
		variations = []string{"METIN:", "METİN:", "METIN :", "METİN :", "CONTENT:", "TEXT:", "Metin:", "Metin :", "Content:", "Text:", "metin:", "metin :", "content:", "text:"}
	} else {
		variations = []string{prefix}
	}
	
	for _, v := range variations {
		idx := strings.Index(text, v)
		if idx != -1 {
			val := text[idx+len(v):]
			return strings.TrimSpace(val)
		}
	}
	
	return ""
}

func truncateContent(content string) string {
	if len(content) > maxContentLength {
		return content[:maxContentLength]
	}
	return content
}

func getVariations(prefix string) []string {
	prefixClean := cleanPrefix(prefix)
	if prefixClean == "BASLIK" {
		return []string{"BASLIK:", "BAŞLIK:", "TITLE:", "Baslik:", "Başlık:", "Title:", "baslik:", "başlık:", "title:"}
	} else if prefixClean == "OZET" {
		return []string{"OZET:", "ÖZET:", "SUMMARY:", "Ozet:", "Özet:", "Summary:", "ozet:", "özet:", "summary:"}
	} else if prefixClean == "METIN" {
		return []string{"METIN:", "METİN:", "METIN :", "METİN :", "CONTENT:", "TEXT:", "Metin:", "Metin :", "Content:", "Text:", "metin:", "metin :", "content:", "text:"}
	}
	return []string{prefix}
}

func extractFieldBetween(text, startPrefix, endPrefix string) string {
	startVariations := getVariations(startPrefix)
	endVariations := getVariations(endPrefix)

	var startIdx = -1
	for _, sv := range startVariations {
		idx := strings.Index(text, sv)
		if idx != -1 {
			startIdx = idx + len(sv)
			break
		}
	}

	if startIdx == -1 {
		return ""
	}

	subText := text[startIdx:]
	var endIdx = -1
	for _, ev := range endVariations {
		idx := strings.Index(subText, ev)
		if idx != -1 {
			endIdx = idx
			break
		}
	}

	if endIdx != -1 {
		return strings.TrimSpace(subText[:endIdx])
	}

	return strings.TrimSpace(subText)
}

