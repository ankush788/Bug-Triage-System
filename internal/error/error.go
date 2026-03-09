package err

import "errors"

// ErrNotFound is returned by repository methods when the requested row does not exist.
var ErrNotFound = errors.New("record not found")

// ErrBugNotFound is returned when a requested bug does not exist.
var ErrBugNotFound = errors.New("bug not found")

// ErrUserNotFound is returned when a lookup by email/ID failed to find a user.
var ErrUserNotFound = errors.New("user not found")