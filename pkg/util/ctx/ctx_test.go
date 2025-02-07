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
			context.WithValue(context.Background(), ctxutil.Key, ctxutil.ValueType{
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

			got := ctxutil.GetFields(tt.ctx)
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
			context.WithValue(context.Background(), ctxutil.Key, ctxutil.ValueType{
				"key": "value",
			}),
		},
		{
			"Overwrite kv",
			context.WithValue(context.Background(), ctxutil.Key, ctxutil.ValueType{
				"key1": "value1",
				"key2": "value2",
			}),
			map[any]any{
				"key1": "aloha",
			},
			context.WithValue(context.Background(), ctxutil.Key, ctxutil.ValueType{
				"key1": "aloha",
				"key2": "value2",
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := ctxutil.PutFields(tt.ctx, tt.fields)
			if tt.want == nil && ctx != nil {
				t.Fatalf("Want ctx to be nil, got %+v", ctx)
			}
			want := ctxutil.GetFields(tt.want)
			got := ctxutil.GetFields(ctx)
			if len(want) == 0 && len(got) == 0 {
				return
			}
			if !reflect.DeepEqual(want, got) {
				t.Fatalf("Want %+v, got %+v", want, got)
			}
		})
	}
}
