package aianalyzer

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"github.com/pgvector/pgvector-go"
)

// TrainAIWithBug embeds a resolved bug report and stores the resulting
// vector in the historical_bugs table.  This is the "seeder" step of
// the RAG workflow.
func (a *GeminiAnalyzer) TrainAIWithBug(ctx context.Context, title, desc, priority, category string) error {
    // generate an embedding for the combined title+description
    res, err := a.embed.EmbedContent(ctx, genai.Text(title+" "+desc))
    if err != nil {
        return err
    }

    err = a.db.Exec(`
        INSERT INTO historical_bugs (title, description, priority, category, embedding)
        VALUES (?, ?, ?, ?, ?)`,
        title, desc, priority, category, pgvector.NewVector(res.Embedding.Values)).Error

    return err
}
