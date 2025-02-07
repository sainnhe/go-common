// Package ctxutil implements utilities for handling context.
package ctxutil

import (
	"context"
)

type keyType struct{}

// Key is the key of the context fields that will be printed in logger.
var Key = keyType{}

// ValueType is the value of the context fields that will be printed in logger.
type ValueType map[any]any

// GetFields get the log fields.
func GetFields(ctx context.Context) map[any]any {
	if ctx == nil {
		return map[any]any{}
	}
	v, ok := ctx.Value(Key).(ValueType)
	if !ok {
		return map[any]any{}
	}
	return map[any]any(v)
}

// PutFields puts the log fields and returns a new context.
func PutFields(ctx context.Context, fields map[any]any) context.Context {
	if ctx == nil {
		return nil
	}
	currentFields := GetFields(ctx)
	if currentFields == nil {
		currentFields = map[any]any{}
	}
	for k, v := range fields {
		currentFields[k] = v
	}
	return context.WithValue(ctx, Key, ValueType(currentFields))
}
