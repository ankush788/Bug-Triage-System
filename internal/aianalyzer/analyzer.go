package aianalyzer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

type GeminiAnalyzer struct {
	model  *genai.GenerativeModel
	logger *zap.Logger
}

type AIResponse struct {
	Priority string `json:"priority"`
	Category string `json:"category"`
}

func NewGeminiAnalyzer(logger *zap.Logger) (*GeminiAnalyzer, error) {
    apiKey := os.Getenv("GEMINI_KEY")
    ctx := context.Background()

    client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
    if err != nil {
        return nil, err
    }

    model := client.GenerativeModel("gemini-3-flash-preview") // Fixed name
    
    // Force the model to output valid JSON
    model.ResponseMIMEType = "application/json"

    return &GeminiAnalyzer{
        model:  model,
        logger: logger,
    }, nil
}

func (a *GeminiAnalyzer) AnalyzeBug(
	ctx context.Context,
	title string,
	description string,
) (string, string, error) {

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	prompt := `
You are a bug triage AI.

Analyze the following software bug and return ONLY valid JSON.

Allowed values:

priority: LOW | MEDIUM | HIGH
category: API | Database | UI | Performance | Security | Other

Return format:

{
 "priority": "HIGH",
 "category": "Database"
}

Bug Title: ` + title + `
Bug Description: ` + description

	resp, err := a.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		a.logger.Error("gemini request failed", zap.Error(err))
		return "", "", err
	}

	if len(resp.Candidates) == 0 {
		return "", "", errors.New("empty response from gemini")
	}
    
	content := resp.Candidates[0].Content.Parts[0].(genai.Text)

	text := string(content)

	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")

	if start == -1 || end == -1 {
		return "", "", errors.New("invalid JSON response from AI")
	}

	cleanJSON := text[start : end+1]

	var result AIResponse

	err = json.Unmarshal([]byte(cleanJSON), &result)
	if err != nil {
		a.logger.Error("failed to parse AI response", zap.Error(err))
		return "", "", err
	}
    
	fmt.Println("This is response ", result)
	a.logger.Debug(
		"bug analyzed using AI",
		zap.String("priority", result.Priority),
		zap.String("category", result.Category),
	)

	return result.Priority, result.Category, nil
}