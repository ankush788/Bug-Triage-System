package aianalyzer

import (
	"context"

	"go.uber.org/zap"
)

// AIResponse represents the priority/category pair returned by the
// generative model.  It is intentionally small to mirror the JSON
// produced by the prompt used in AnalyzeBug.
type AIResponse struct {
    Priority string `json:"priority"`
    Category string `json:"category"`
}

// Analyzer defines the interface for AI-based bug analysis.
// Any analyzer implementation must provide methods to analyze bugs and train with new data.
type Analyzer interface {
    AnalyzeBug(logger *zap.Logger, ctx context.Context, title, description string) (string, string, error)
    TrainAIWithBug(ctx context.Context, title, desc, priority, category string) error
    getContext(logger *zap.Logger, ctx context.Context, title, description string) (string, error)   
    
}
