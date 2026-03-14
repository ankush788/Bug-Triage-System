package aianalyzer

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/pgvector/pgvector-go"
	"go.uber.org/zap"
)

// getContext performs a vector similarity search against the
// historical_bugs table and returns a back‑quoted list of the top
// examples.  It is intentionally unexported; callers should use it via
// AnalyzeBug.
func (a *GeminiAnalyzer) getContext(logger *zap.Logger, ctx context.Context, title, description string) (string, error) {
    logger.Debug("Fetching historical context", zap.String("title", title), zap.String("description", description))
    res, err := a.embed.EmbedContent(ctx, genai.Text(title+" "+description))
    if err != nil {
        logger.Error("Failed to embed content", zap.Error(err))
        return "", err
    }

    //2. Vector Similarity Operators
    //   Euclidean distance
    //   Inner product
    //   Cosine distance 

    //Take stored embeddings from historical_bugs, compare them with the new bug embedding, sort by similarity,
    //  and return the 3 most similar bugs. and give it's title , priority , cateogry
    
    rows, err := a.db.Raw(`
        SELECT title, priority, category 
        FROM historical_bugs 
        ORDER BY embedding <=> ? 
        LIMIT 3`, pgvector.NewVector(res.Embedding.Values)).Rows()
    if err != nil {
        logger.Error("Failed to query historical bugs", zap.Error(err))
        return "", err
    }
    defer rows.Close()
    
    var examples []string
    for rows.Next() {
        var t, p, c string
        if err := rows.Scan(&t, &p, &c); err == nil {
            examples = append(examples, fmt.Sprintf("Past Bug: %s -> Label: [%s, %s]", t, p, c))
        }
    }

    if len(examples) == 0 {
        logger.Debug("No historical data found for the given bug report")
        return "No historical data found.", nil
    }

    return strings.Join(examples, "\n"), nil
}
