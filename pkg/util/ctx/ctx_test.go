package ctxutil_test

import (
	"context"
	"reflect"
	"testing"

	ctxutil "github.com/teamsorghum/go-common/pkg/util/ctx"
)

func TestContext_GetContextFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		ctx            context.Context
		expectedFields map[any]any
	}{
		{
			"EmptyFields",
			context.Background(),
			map[any]any{},
		},
		{
			"NonEmptyFields",
			context.WithValue(context.Background(), ctxutil.ContextKey, ctxutil.ContextValueType{
				"key": "value",
			}),
			map[any]any{
				"key": "value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fields := ctxutil.GetContextFields(tt.ctx)
			if !reflect.DeepEqual(fields, tt.expectedFields) {
				t.Fatalf("expectedFields = %+v, get %+v", tt.expectedFields, fields)
			}
		})
	}
}

func TestContext_PutContextFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		ctx         context.Context
		fields      map[any]any
		expectedCtx context.Context
	}{
		{
			"NilContext",
			nil,
			map[any]any{},
			nil,
		},
		{
			"NilFields",
			context.Background(),
			nil,
			context.Background(),
		},
		{
			"AddKV",
			context.Background(),
			map[any]any{
				"key": "value",
			},
			context.WithValue(context.Background(), ctxutil.ContextKey, ctxutil.ContextValueType{
				"key": "value",
			}),
		},
		{
			"OverwriteKV",
			context.WithValue(context.Background(), ctxutil.ContextKey, ctxutil.ContextValueType{
				"key1": "value1",
				"key2": "value2",
			}),
			map[any]any{
				"key1": "aloha",
			},
			context.WithValue(context.Background(), ctxutil.ContextKey, ctxutil.ContextValueType{
				"key1": "aloha",
				"key2": "value2",
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := ctxutil.PutContextFields(tt.ctx, tt.fields)
			if tt.expectedCtx == nil && ctx != nil {
				t.Fatalf("Expect ctx to be nil, get %+v", ctx)
			}
			expectedFields := ctxutil.GetContextFields(tt.expectedCtx)
			actualFields := ctxutil.GetContextFields(ctx)
			if len(expectedFields) == 0 && len(actualFields) == 0 {
				return
			}
			if !reflect.DeepEqual(expectedFields, actualFields) {
				t.Fatalf("expectedCtx == %+v, actualFields == %+v", expectedFields, actualFields)
			}
		})
	}
}
