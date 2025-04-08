package otel_test

import (
	"context"
	"testing"

	"github.com/teamsorghum/go-common/pkg/encoding"
	"github.com/teamsorghum/go-common/pkg/log"
	"github.com/teamsorghum/go-common/pkg/otel"
	gotel "go.opentelemetry.io/otel"
)

func TestNew_disabled(t *testing.T) {
	t.Parallel()

	p, tp, mp, lp, c, err := otel.New(&otel.Config{Enable: false})
	if p == nil || tp == nil || mp == nil || lp == nil || c == nil || err != nil {
		t.Fatal("Expect err == nil, and other returns != nil.")
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		getConfig       func() *otel.Config
		expectInitError bool
	}{
		{
			"Nil config",
			func() *otel.Config { return nil },
			true,
		},
		{
			"TLS",
			func() *otel.Config {
				cfg, err := encoding.LoadConfig[otel.Config](nil, encoding.TypeNil)
				if err != nil {
					t.Fatal(err.Error())
				}
				cfg.Conn.EnableTLS = true
				cfg.TimeoutMs = 100
				return cfg
			},
			false,
		},
		{
			"Attributes",
			func() *otel.Config {
				cfg, err := encoding.LoadConfig[otel.Config](nil, encoding.TypeNil)
				if err != nil {
					t.Fatal(err.Error())
				}
				cfg.Attributes = map[string]string{
					"attr1": "value1",
					"attr2": "value2",
				}
				return cfg
			},
			false,
		},
		{
			"Headers",
			func() *otel.Config {
				cfg, err := encoding.LoadConfig[otel.Config](nil, encoding.TypeNil)
				if err != nil {
					t.Fatal(err.Error())
				}
				cfg.Headers = map[string]string{
					"header1": "value1",
					"header2": "value2",
				}
				return cfg
			},
			false,
		},
		{
			"Simple processor",
			func() *otel.Config {
				cfg, err := encoding.LoadConfig[otel.Config](nil, encoding.TypeNil)
				if err != nil {
					t.Fatal(err.Error())
				}
				cfg.Batch.MaxSize = 0
				return cfg
			},
			false,
		},
		{
			"Invalid gRPC address",
			func() *otel.Config {
				cfg, err := encoding.LoadConfig[otel.Config](nil, encoding.TypeNil)
				if err != nil {
					t.Fatal(err.Error())
				}
				cfg.Conn.Port = 22
				cfg.TimeoutMs = 100
				return cfg
			},
			false,
		},
		{
			"Cumulative metric temporality",
			func() *otel.Config {
				cfg, err := encoding.LoadConfig[otel.Config](nil, encoding.TypeNil)
				if err != nil {
					t.Fatal(err.Error())
				}
				cfg.Metric.Temporality = "cumulative"
				return cfg
			},
			false,
		},
		{
			"Delta metric temporality",
			func() *otel.Config {
				cfg, err := encoding.LoadConfig[otel.Config](nil, encoding.TypeNil)
				if err != nil {
					t.Fatal(err.Error())
				}
				cfg.Metric.Temporality = "delta"
				return cfg
			},
			false,
		},
		{
			"Invalid metric temporality",
			func() *otel.Config {
				cfg, err := encoding.LoadConfig[otel.Config](nil, encoding.TypeNil)
				if err != nil {
					t.Fatal(err.Error())
				}
				cfg.Metric.Temporality = "nil"
				return cfg
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p, tp, mp, lp, c, err := otel.New(tt.getConfig())
			if tt.expectInitError {
				if err == nil {
					t.Fatal("Expect error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("%+v", err)
			}
			if p == nil || tp == nil || mp == nil || lp == nil || c == nil {
				t.Fatalf("Expect non nil, got p = %+v, tp = %+v, mp = %+v, lp = %+v", p, tp, mp, lp)
			}
			defer c()

			tracer := gotel.Tracer("test")
			meter := gotel.Meter("test")
			logger := log.WithOTelAttrs(log.GetGlobalLogger())

			// Start tracer
			ctx, span := tracer.Start(context.Background(), "test")
			defer span.End()

			// Initialize metric counter
			counter, err := meter.Int64Counter("test")
			if err != nil {
				t.Fatal(err.Error())
			}

			// Increase counter and print a log
			counter.Add(ctx, 1)
			logger.InfoContext(ctx, "Hello world!")
		})
	}
}
