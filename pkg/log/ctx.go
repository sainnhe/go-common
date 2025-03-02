package log

import (
	"context"
	"maps"
)

type ctxKeyT struct{}

// ctxKey is the key of the context fields that will be tracked by logger.
var ctxKey = ctxKeyT{}

// ctxValueT is the value type of the context fields that will be tracked by logger.
type ctxValueT map[any]any

// GetCtxFields extracts log fields from a context.
func GetCtxFields(ctx context.Context) map[any]any {
	if ctx == nil {
		return map[any]any{}
	}
	v, ok := ctx.Value(ctxKey).(ctxValueT)
	if !ok {
		return map[any]any{}
	}
	return map[any]any(v)
}

// PutCtxFields puts log fields into a new context and return it.
func PutCtxFields(ctx context.Context, fields map[any]any) context.Context {
	if ctx == nil {
		return nil
	}
	currentFields := GetCtxFields(ctx)
	maps.Copy(currentFields, fields)
	return context.WithValue(ctx, ctxKey, ctxValueT(currentFields))
}
