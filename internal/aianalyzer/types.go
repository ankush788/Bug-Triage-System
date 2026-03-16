package aianalyzer

import (
	"context"

	"go.uber.org/zap"
)


// Analyzer defines the interface for AI-based bug analysis.
type Analyzer interface {
    AnalyzeBug(logger *zap.Logger, ctx context.Context, title, description string) (string, string, error)
    TrainAIWithBug(ctx context.Context, title, desc, priority, category string) error
    GetContext(logger *zap.Logger, ctx context.Context, title, description string) (string, error)   
    
}


type AIResponse struct {
    Priority string `json:"priority"`
    Category string `json:"category"`
}

