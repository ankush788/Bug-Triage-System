package aianalyzer

import (
	"context"
	"math/rand"
	"strings"
	"go.uber.org/zap"
)

// SimpleAIAnalyzer performs basic heuristic-based bug classification
// In production, this would call OpenAI or Gemini API
type SimpleAIAnalyzer struct {
	logger *zap.Logger
}

func NewSimpleAIAnalyzer(logger *zap.Logger) *SimpleAIAnalyzer {
	return &SimpleAIAnalyzer{logger: logger}
}

// AnalyzeBug analyzes bug title and description to assign priority and category
func (a *SimpleAIAnalyzer) AnalyzeBug(ctx context.Context, title, description string) (priority, category string, err error) {
	combined := strings.ToLower(title + " " + description)

	// Simple heuristic-based priority assignment
	priority = "MEDIUM"

	criticalKeywords := []string{"crash", "fatal", "critical", "urgent", "broken", "cannot", "unable"}
	for _, keyword := range criticalKeywords {
		if strings.Contains(combined, keyword) {
			priority = "HIGH"
			break
		}
	}

	minorKeywords := []string{"typo", "formatting", "minor", "cosmetic", "ui"}
	for _, keyword := range minorKeywords {
		if strings.Contains(combined, keyword) {
			priority = "LOW"
			break
		}
	}

	// Simple category assignment
	categories := []string{"API", "Database", "UI", "Performance", "Security", "Other"}
	category = categories[rand.Intn(len(categories))]

	if strings.Contains(combined, "database") || strings.Contains(combined, "query") {
		category = "Database"
	} else if strings.Contains(combined, "api") || strings.Contains(combined, "endpoint") {
		category = "API"
	} else if strings.Contains(combined, "ui") || strings.Contains(combined, "button") || strings.Contains(combined, "display") {
		category = "UI"
	} else if strings.Contains(combined, "slow") || strings.Contains(combined, "performance") || strings.Contains(combined, "timeout") {
		category = "Performance"
	} else if strings.Contains(combined, "security") || strings.Contains(combined, "password") || strings.Contains(combined, "token") {
		category = "Security"
	}

	a.logger.Debug("bug analyzed",
		zap.String("priority", priority),
		zap.String("category", category),
	)

	return priority, category, nil
}
