package log

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/grafana/loki-client-go/loki"
	slogloki "github.com/samber/slog-loki/v3"
	slogotel "github.com/samber/slog-otel"
)

// NewLoki initializes a Loki based logger.
func NewLoki(cfg *Config) (logger Logger, cleanup func(), err error) {
	if cfg == nil {
		err = errors.New("nil dependency")
		return
	}
	// Setup loki client
	lokiCfg, err := loki.NewDefaultConfig(cfg.Loki.URL)
	if err != nil {
		return
	}
	lokiCfg.TenantID = cfg.Loki.TenantID
	lokiCfg.Timeout = time.Duration(cfg.Loki.TimeoutSec) * time.Second
	if err = json.Unmarshal([]byte(cfg.Loki.ExternalLabels), &lokiCfg.ExternalLabels); err != nil {
		return
	}
	client, err := loki.New(lokiCfg)
	if err != nil {
		return
	}
	cleanup = func() {
		client.Stop()
	}
	// Setup slog
	var logLevel slog.Level
	switch cfg.Level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		return nil, nil, errors.New("invalid log level")
	}
	slogLogger := slog.New(slogloki.Option{
		Level:     logLevel,
		Client:    client,
		AddSource: true,
		AttrFromContext: []func(ctx context.Context) []slog.Attr{
			slogotel.ExtractOtelAttrFromContext([]string{"trace"}, "trace_id", "span_id"),
		},
	}.NewLokiHandler())
	logger = &slogImpl{
		slogLogger,
		slogLogger,
		nil,
		[]any{},
	}
	return
}
