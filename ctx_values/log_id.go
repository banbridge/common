package ctx_values

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

// logIDContextKey is how we find logID in a context.Context.
// Unexported to prevent context key collisions.
type logIDContextKey string

var (
	logIDKey      = logIDContextKey("logID")
	nonceLogIDKey = logIDContextKey("nonceLogID")
)

// NewContextWithLogID returns a context derived from ctx that carries logID.
func NewContextWithLogID(ctx context.Context, logID string) context.Context {
	return context.WithValue(ctx, logIDKey, logID)
}

// LogIDFromContext returns the log ID from ctx, or "-" when ctx does not carry one.
func LogIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(logIDKey).(string); ok {
		return v
	}
	return "-"
}

// NewContextWithNonceLogID returns a context derived from ctx that carries nonceLogID.
func NewContextWithNonceLogID(ctx context.Context, nonceLogID string) context.Context {
	return context.WithValue(ctx, nonceLogIDKey, nonceLogID)
}

// NonceLogIDFromContext returns the nonce log ID from ctx, or "-" when ctx does not carry one.
func NonceLogIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(nonceLogIDKey).(string); ok {
		return v
	}
	return "-"
}

// NewContextWithNamedLogID returns a context derived from ctx that carries logID under name.
// The name argument is typically an internal key to avoid context key collisions.
func NewContextWithNamedLogID(ctx context.Context, name any, logID string) context.Context {
	return context.WithValue(ctx, name, logID)
}

// NamedLogIDFromContext returns the named log ID from ctx, or "-" when ctx does not carry one.
func NamedLogIDFromContext(ctx context.Context, name any) string {
	if v, ok := ctx.Value(name).(string); ok {
		return v
	}
	return "-"
}

// NewShortLogID creates a short log ID using the final UUID segment.
func NewShortLogID() string {
	items := strings.Split(uuid.New().String(), "-")
	return items[len(items)-1]
}
