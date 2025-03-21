package otel_test

import (
	"context"
	"fmt"

	"github.com/teamsorghum/go-common/pkg/constant"
	"github.com/teamsorghum/go-common/pkg/encoding"
	"github.com/teamsorghum/go-common/pkg/log"
	"github.com/teamsorghum/go-common/pkg/otel"
	gotel "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// This example demonstrates how to use this package to initialize and use propagator and providers.
func Example_usage() {
	// Initialize config.
	cfg, err := encoding.LoadConfig[otel.Config](nil, encoding.TypeNil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Modify some options to make it more suitable for testing environment.
	cfg.Batch.MaxDelayMs = 1000
	cfg.Trace.AlwaysSample = true
	cfg.Metric.ReaderIntervalMs = 3000

	// Instantiate new propagator and providers and set them as global.
	// The first 4 returned values are propagator and providers.
	// Since they are already set as global, we ignore them here.
	_, _, _, _, cleanup, err := otel.New(cfg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer cleanup()

	// Initialize a context with baggage.
	// The information contained in this baggage will be spread across servers carried by this context.
	property, err := baggage.NewKeyValueProperty("property_key", "property_value")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	member, err := baggage.NewMember("member_key", "member_value", property)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	b, err := baggage.New(member)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ctx := baggage.ContextWithBaggage(context.Background(), b)

	// Attributes represent additional key-value descriptors that can be bound to a metric observer or recorder.
	attributes := []attribute.KeyValue{
		attribute.String("attrA", "chocolate"),
		attribute.String("attrB", "raspberry"),
		attribute.String("attrC", "vanilla"),
	}

	// We can extract baggage information from context and append them to attributes.
	members := baggage.FromContext(ctx).Members()
	for _, member := range members {
		attributes = append(attributes, attribute.KeyValue{
			Key:   attribute.Key(member.Key()),
			Value: attribute.StringValue(member.Value()),
		})
		for _, property := range member.Properties() {
			value, ok := property.Value()
			if !ok {
				value = "nil"
			}
			attributes = append(attributes, attribute.KeyValue{
				Key:   attribute.Key(property.Key()),
				Value: attribute.StringValue(value),
			})
		}
	}

	// Initialize tracer, meter and logger
	name := "github.com/teamsorghum/go-common/pkg/otel"
	tracer := gotel.Tracer(name)
	meter := gotel.Meter(name, metric.WithInstrumentationAttributes(attributes...))
	logger := log.WithOTelAttrs(log.NewOTel(name), attributes...)

	// Start tracer
	ctx, span := tracer.Start(ctx, "Example", trace.WithAttributes(attributes...))
	defer span.End()

	// Initialize metric counter
	counter, err := meter.Int64Counter(name)
	if err != nil {
		logger.ErrorContext(ctx, "Init int64 counter error.", constant.LogAttrError, err)
		return
	}

	// Increase counter and print a log
	counter.Add(ctx, 1)
	logger.InfoContext(ctx, "Hello world!")

	// Output:
}
