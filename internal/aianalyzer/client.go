package aianalyzer

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"go.uber.org/zap"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

// GeminiAnalyzer wraps the generative and embedding models along with
// database and logging dependencies.
type GeminiAnalyzer struct {
    model  *genai.GenerativeModel
    embed  *genai.EmbeddingModel
    db     *gorm.DB
    logger *zap.Logger
}

// NewGeminiAnalyzer constructs an analyzer instance.  A working
// *gorm.DB must be supplied.
// The GEMINI_KEY environment variable is used to authenticate the AI
// client.
func NewGeminiAnalyzer(logger *zap.Logger, db *gorm.DB) (*GeminiAnalyzer, error) {
    // apiKey := os.Getenv("GEMINI_KEY")
    apiKey := os.Getenv("GEMINI_KEY1")
    geminiModel := os.Getenv("GEMINI_MODEL")

    // Safety check: Don't let the app start with empty credentials
    if apiKey == "" || geminiModel == "" {
        return nil, fmt.Errorf("GEMINI_KEY or GEMINI_MODEL not set in environment")
    }
    
    ctx := context.Background()
    client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
    if err != nil {
        return nil, err
    }

    model := client.GenerativeModel(geminiModel)
    embedModel := client.EmbeddingModel("gemini-embedding-001")

    // make sure responses come back as JSON so we can unmarshal them
    model.ResponseMIMEType = "application/json"

    return &GeminiAnalyzer{
        model:  model,
        embed:  embedModel,
        db:     db,
        logger: logger,
    }, nil
}
