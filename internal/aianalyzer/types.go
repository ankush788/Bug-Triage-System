package aianalyzer

// AIResponse represents the priority/category pair returned by the
// generative model.  It is intentionally small to mirror the JSON
// produced by the prompt used in AnalyzeBug.
type AIResponse struct {
    Priority string `json:"priority"`
    Category string `json:"category"`
}
