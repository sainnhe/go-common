// Package otel implements [OpenTelemetry] related functions.
//
// [OpenTelemetry]: https://opentelemetry.io/
package otel

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"time"

	"github.com/teamsorghum/go-common/pkg/constant"
	clog "github.com/teamsorghum/go-common/pkg/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
)

const compressor = "gzip"

// ErrInvalidConfig indicates the given config is invalid.
var ErrInvalidConfig = errors.New("invalid config")

// New instantiates a new [propagation.TextMapPropagator], [trace.TracerProvider], [metric.MeterProvider] and
// [log.LoggerProvider], and sets them as the global propagator and providers.
//
// NOTE: The returned cleanup function will handle shutdown correctly, so you don't need to manually call the shutdown
// functions of returned providers.
func New(cfg *Config) (propagator propagation.TextMapPropagator, tracerProvider *trace.TracerProvider,
	meterProvider *metric.MeterProvider, loggerProvider *log.LoggerProvider, cleanup func(), err error) {
	// Check argument
	if cfg == nil {
		err = constant.ErrNilDeps
		return
	}

	// Return if disabled
	if !cfg.Enable {
		propagator = propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		)
		tracerProvider = trace.NewTracerProvider()
		meterProvider = metric.NewMeterProvider()
		loggerProvider = log.NewLoggerProvider()
		cleanup = func() {}
		return
	}

	// Base endpoint URL
	baseEndpointURL := ""
	if cfg.Conn.EnableTLS {
		baseEndpointURL = fmt.Sprintf("https://%s:%d", cfg.Conn.Host, cfg.Conn.Port)
	} else {
		baseEndpointURL = fmt.Sprintf("http://%s:%d", cfg.Conn.Host, cfg.Conn.Port)
	}

	// Credentials
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return
	}
	creds := credentials.NewTLS(&tls.Config{
		RootCAs:    rootCAs,
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,
	})

	// Timeout
	timeout := time.Duration(cfg.TimeoutMs) * time.Millisecond

	// Context
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Resource
	attrs := make([]attribute.KeyValue, 0, len(cfg.Attributes))
	for k, v := range cfg.Attributes {
		attrs = append(attrs, attribute.KeyValue{
			Key:   attribute.Key(k),
			Value: attribute.StringValue(v),
		})
	}
	res, err := resource.New(ctx,
		resource.WithAttributes(attrs...),
		resource.WithContainer(),
		resource.WithFromEnv(),
		resource.WithHost(),
		resource.WithOS(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		return
	}

	// Propagator
	propagator = propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	// Tracer provider
	tracerProvider, err = initTracerProvider(
		ctx, cfg, baseEndpointURL+cfg.Trace.Path, timeout, creds, res,
	)
	if err != nil {
		return
	}

	// Meter provider
	meterProvider, err = initMeterProvider(
		ctx, cfg, baseEndpointURL+cfg.Metric.Path, timeout, creds, res,
	)
	if err != nil {
		return
	}

	// Logger provider
	loggerProvider, err = initLoggerProvider(
		ctx, cfg, baseEndpointURL+cfg.Log.Path, timeout, creds, res,
	)
	if err != nil {
		return
	}

	// Set as global propagator and providers.
	otel.SetTextMapPropagator(propagator)
	otel.SetTracerProvider(tracerProvider)
	otel.SetMeterProvider(meterProvider)
	global.SetLoggerProvider(loggerProvider)

	// Cleanup
	cleanup = func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if err := errors.Join(
			tracerProvider.ForceFlush(ctx),
			tracerProvider.Shutdown(ctx),
			meterProvider.ForceFlush(ctx),
			meterProvider.Shutdown(ctx),
			loggerProvider.ForceFlush(ctx),
			loggerProvider.Shutdown(ctx),
		); err != nil {
			clog.NewLogger("github.com/teamsorghum/go-common/pkg/otel").ErrorContext(
				ctx, "Cleanup error.", constant.LogAttrError, err)
		}
	}

	return
}

func initTracerProvider(
	ctx context.Context, cfg *Config, endpointURL string, timeout time.Duration, creds credentials.TransportCredentials,
	res *resource.Resource) (provider *trace.TracerProvider, err error) {
	// Exporter
	exporterOpts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpointURL(endpointURL),
		otlptracegrpc.WithTimeout(timeout),
	}
	if cfg.Conn.EnableTLS {
		exporterOpts = append(exporterOpts, otlptracegrpc.WithTLSCredentials(creds))
	} else {
		exporterOpts = append(exporterOpts, otlptracegrpc.WithInsecure())
	}
	if cfg.EnableGzip {
		exporterOpts = append(exporterOpts, otlptracegrpc.WithCompressor(compressor))
	}
	if len(cfg.Headers) > 0 {
		exporterOpts = append(exporterOpts, otlptracegrpc.WithHeaders(cfg.Headers))
	}
	exporter, err := otlptracegrpc.New(ctx, exporterOpts...)
	if err != nil {
		return
	}

	// Provider
	providerOpts := []trace.TracerProviderOption{
		trace.WithResource(res),
	}
	if cfg.Trace.AlwaysSample {
		providerOpts = append(providerOpts, trace.WithSampler(trace.AlwaysSample()))
	}
	if cfg.Batch.MaxSize > 0 {
		providerOpts = append(providerOpts, trace.WithSpanProcessor(trace.NewBatchSpanProcessor(exporter,
			trace.WithMaxExportBatchSize(cfg.Batch.MaxSize),
			trace.WithMaxQueueSize(cfg.Batch.QueueSize),
			trace.WithBatchTimeout(time.Duration(cfg.Batch.MaxDelayMs)*time.Millisecond),
			trace.WithExportTimeout(timeout),
		)))
	} else {
		providerOpts = append(providerOpts, trace.WithSpanProcessor(trace.NewSimpleSpanProcessor(exporter)))
	}
	provider = trace.NewTracerProvider(providerOpts...)

	return
}

func initMeterProvider(
	ctx context.Context, cfg *Config, endpointURL string, timeout time.Duration, creds credentials.TransportCredentials,
	res *resource.Resource) (provider *metric.MeterProvider, err error) {
	// Exporter
	exporterOpts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpointURL(endpointURL),
		otlpmetricgrpc.WithTimeout(timeout),
	}
	switch cfg.Metric.Temporality {
	case "default":
		break
	case "cumulative":
		exporterOpts = append(exporterOpts, otlpmetricgrpc.WithTemporalitySelector(
			func(_ metric.InstrumentKind) metricdata.Temporality {
				return metricdata.CumulativeTemporality
			}))
	case "delta":
		exporterOpts = append(exporterOpts, otlpmetricgrpc.WithTemporalitySelector(
			func(_ metric.InstrumentKind) metricdata.Temporality {
				return metricdata.DeltaTemporality
			}))
	default:
		err = ErrInvalidConfig
		return
	}
	if cfg.Conn.EnableTLS {
		exporterOpts = append(exporterOpts, otlpmetricgrpc.WithTLSCredentials(creds))
	} else {
		exporterOpts = append(exporterOpts, otlpmetricgrpc.WithInsecure())
	}
	if cfg.EnableGzip {
		exporterOpts = append(exporterOpts, otlpmetricgrpc.WithCompressor(compressor))
	}
	if len(cfg.Headers) > 0 {
		exporterOpts = append(exporterOpts, otlpmetricgrpc.WithHeaders(cfg.Headers))
	}
	exporter, err := otlpmetricgrpc.New(ctx, exporterOpts...)
	if err != nil {
		return
	}

	// Provider
	providerOpts := []metric.Option{
		metric.WithResource(res),
		metric.WithReader(
			metric.NewPeriodicReader(
				exporter, metric.WithInterval(time.Duration(cfg.Metric.ReaderIntervalMs)*time.Millisecond),
			),
		),
	}
	provider = metric.NewMeterProvider(providerOpts...)

	return
}

func initLoggerProvider(
	ctx context.Context, cfg *Config, endpointURL string, timeout time.Duration, creds credentials.TransportCredentials,
	res *resource.Resource) (provider *log.LoggerProvider, err error) {
	// Exporter
	exporterOpts := []otlploggrpc.Option{
		otlploggrpc.WithEndpointURL(endpointURL),
		otlploggrpc.WithTimeout(timeout),
	}
	if cfg.Conn.EnableTLS {
		exporterOpts = append(exporterOpts, otlploggrpc.WithTLSCredentials(creds))
	} else {
		exporterOpts = append(exporterOpts, otlploggrpc.WithInsecure())
	}
	if cfg.EnableGzip {
		exporterOpts = append(exporterOpts, otlploggrpc.WithCompressor(compressor))
	}
	if len(cfg.Headers) > 0 {
		exporterOpts = append(exporterOpts, otlploggrpc.WithHeaders(cfg.Headers))
	}
	exporter, err := otlploggrpc.New(ctx, exporterOpts...)
	if err != nil {
		return
	}

	// Provider
	providerOpts := []log.LoggerProviderOption{
		log.WithResource(res),
	}
	if cfg.Batch.MaxSize > 0 {
		providerOpts = append(providerOpts, log.WithProcessor(log.NewBatchProcessor(exporter,
			log.WithExportMaxBatchSize(cfg.Batch.MaxSize),
			log.WithMaxQueueSize(cfg.Batch.QueueSize),
			log.WithExportInterval(time.Duration(cfg.Batch.MaxDelayMs)*time.Millisecond),
			log.WithExportTimeout(timeout),
		)))
	} else {
		providerOpts = append(providerOpts, log.WithProcessor(log.NewSimpleProcessor(exporter)))
	}
	provider = log.NewLoggerProvider(providerOpts...)

	return
}
