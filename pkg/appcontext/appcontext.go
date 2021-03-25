package appcontext

import (
	"context"
)

type contextKey string

const (
	// KeyUserID represents the current logged-in UserID
	KeyUserID contextKey = "UserID"

	// KeySessionID represents the current logged-in SessionID
	KeySessionID contextKey = "SessionID"
)

// UserID gets current userId logged in from the context
func UserID(ctx context.Context) int {
	userID := (ctx).Value(KeyUserID)
	if userID != nil {
		v := userID.(int)
		return v
	}
	return 0
}

// UserID gets current userId logged in from the context
func SessionID(ctx context.Context) string {
	sessionID := (ctx).Value(KeySessionID)
	if sessionID != nil {
		v := sessionID.(string)
		return v
	}
	return ""
}
