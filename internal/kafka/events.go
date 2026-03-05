package kafka

import "encoding/json"

// Event types
const (
	EventBugCreated  = "bug_created"
	EventBugAnalyzed = "bug_analyzed"
)

// BugCreatedEvent published when a new bug is created
type BugCreatedEvent struct {
	BugID       int64  `json:"bug_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ReporterID  int64  `json:"reporter_id"`
}

// ToJSON converts event to JSON bytes
func (e *BugCreatedEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// ParseBugCreatedEvent parses JSON into BugCreatedEvent
func ParseBugCreatedEvent(data []byte) (*BugCreatedEvent, error) {
	var event BugCreatedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

// BugAnalyzedEvent published after AI analysis completes
type BugAnalyzedEvent struct {
	BugID    int64  `json:"bug_id"`
	Priority string `json:"priority"`
	Category string `json:"category"`
}

// ToJSON converts event to JSON bytes
func (e *BugAnalyzedEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// ParseBugAnalyzedEvent parses JSON into BugAnalyzedEvent
func ParseBugAnalyzedEvent(data []byte) (*BugAnalyzedEvent, error) {
	var event BugAnalyzedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}
