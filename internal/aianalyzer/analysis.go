package aianalyzer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/generative-ai-go/genai"
	"go.uber.org/zap"
)

// AnalyzeBug :- augments a new bug report with examples from history and asks the generative model to assign a priority and category.  It
// returns the two labels (priority, category) or an error if the request fails.
func (a *GeminiAnalyzer) AnalyzeBug(logger *zap.Logger, ctx context.Context, title, description string) (string, string, error) {
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()

    pastContext, err := a.getContext(logger ,ctx, title, description)
    // default prompt (for  trying better search results when prompt is empty)
    if err != nil || pastContext == "" {
        pastContext = "No historical examples available. Classify based on bug severity and type."
    }

    prompt := fmt.Sprintf(`
You are a bug triage AI. Categorize the NEW BUG based on the provided EXAMPLES from our history. 

### EXAMPLES:
%s

### NEW BUG:
Title: %s
Description: %s

Return ONLY JSON: {"priority": "...", "category": "..."}
`, pastContext, title, description)

    resp, err := a.model.GenerateContent(ctx, genai.Text(prompt))
    if err != nil {
        return "", "", err
    }
    
    if len(resp.Candidates) == 0 {
	return "", "", fmt.Errorf("no response from model")
    }

    content := resp.Candidates[0].Content.Parts[0].(genai.Text)
    var result AIResponse
    if err := json.Unmarshal([]byte(string(content)), &result); err != nil {
        return "", "", err
    }

    // use this bug report as a new training example for the future
    if err := a.TrainAIWithBug(ctx, title, description, result.Priority, result.Category) ;err != nil {
        a.logger.Error("Failed to train AI with new bug report", zap.Error(err))
    }
    return result.Priority, result.Category, nil
}
