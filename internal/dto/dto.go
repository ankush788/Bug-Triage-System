package dto

// Auth DTOs

// RegisterRequest holds incoming registration data
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterResponse returns user info and auth token
type RegisterResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

// LoginRequest holds incoming login data
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse returns user info and auth token
type LoginResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

// Bug DTOs

// CreateBugRequest holds incoming bug creation data
type CreateBugRequest struct {
	Title       string `json:"title" binding:"required,min=5"`
	Description string `json:"description" binding:"required,min=10"`
}

// UpdateBugStatusRequest holds status update data
type UpdateBugStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// BugResponse represents a bug in API responses
type BugResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	Category    string `json:"category"`
	ReporterID  int64  `json:"reporter_id"`
	CreatedAt   string `json:"created_at,omitempty"` //If this variable is empty/zero, just leave it out of the final JSON result
}

// BugsListResponse represents a paginated list of bugs
type BugsListResponse struct {
	Bugs   []BugResponse `json:"bugs"`
	Limit  int           `json:"limit"`
	Offset int           `json:"offset"`
}