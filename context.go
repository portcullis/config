package config

import (
	"context"
)

type contextKey string

var (
	configSetContextKey = contextKey("config-set")
)

// FromContext extracts the config.Set instance if it exists from the provided context or nil if not found
func FromContext(ctx context.Context) *Set {
	set := ctx.Value(configSetContextKey)
	if set == nil {
		return nil
	}

	return set.(*Set)
}

// NewContext creates a child context of the supplied context embedding the *config.Set. This *config.Set can be retrieved with the FromContext
func NewContext(ctx context.Context, set *Set) context.Context {
	return context.WithValue(ctx, configSetContextKey, set)
}
