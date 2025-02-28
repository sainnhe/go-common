package log_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/teamsorghum/go-common/pkg/log"
)

func TestContext_GetContextFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ctx  context.Context
		want map[any]any
	}{
		{
			"Empty fields",
			context.Background(),
			map[any]any{},
		},
		{
			"Non empty fields",
			context.WithValue(context.Background(), log.ContextKey, log.ContextValueType{
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

			got := log.GetContextFields(tt.ctx)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Want %+v, got %+v", tt.want, got)
			}
		})
	}
}

func TestContext_PutContextFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		ctx    context.Context
		fields map[any]any
		want   context.Context
	}{
		{
			"Nil context",
			nil,
			map[any]any{},
			nil,
		},
		{
			"Nil fields",
			context.Background(),
			nil,
			context.Background(),
		},
		{
			"Add kv",
			context.Background(),
			map[any]any{
				"key": "value",
			},
			context.WithValue(context.Background(), log.ContextKey, log.ContextValueType{
				"key": "value",
			}),
		},
		{
			"Overwrite kv",
			context.WithValue(context.Background(), log.ContextKey, log.ContextValueType{
				"key1": "value1",
				"key2": "value2",
			}),
			map[any]any{
				"key1": "aloha",
			},
			context.WithValue(context.Background(), log.ContextKey, log.ContextValueType{
				"key1": "aloha",
				"key2": "value2",
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := log.PutContextFields(tt.ctx, tt.fields)
			if tt.want == nil && ctx != nil {
				t.Fatalf("Want ctx to be nil, got %+v", ctx)
			}
			want := log.GetContextFields(tt.want)
			got := log.GetContextFields(ctx)
			if len(want) == 0 && len(got) == 0 {
				return
			}
			if !reflect.DeepEqual(want, got) {
				t.Fatalf("Want %+v, got %+v", want, got)
			}
		})
	}
}
