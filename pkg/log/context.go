package log

import "context"

type keyType struct{}

// ContextKey is the key of the context fields that will be printed in logger.
var ContextKey = keyType{}

// ContextValueType is the value of the context fields that will be printed in logger.
type ContextValueType map[any]any

// GetContextFields get the log fields from a context.
func GetContextFields(ctx context.Context) map[any]any {
	if ctx == nil {
		return map[any]any{}
	}
	v, ok := ctx.Value(ContextKey).(ContextValueType)
	if !ok {
		return map[any]any{}
	}
	return map[any]any(v)
}

// PutContextFields puts the log fields and returns a new context.
func PutContextFields(ctx context.Context, fields map[any]any) context.Context {
	if ctx == nil {
		return nil
	}
	currentFields := GetContextFields(ctx)
	if currentFields == nil {
		currentFields = map[any]any{}
	}
	for k, v := range fields {
		currentFields[k] = v
	}
	return context.WithValue(ctx, ContextKey, ContextValueType(currentFields))
}
