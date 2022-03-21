package opentelemetry

import (
	"runtime/debug"

	"github.com/pepol/databuddy/internal/log"
	"go.opentelemetry.io/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

// Tracer returns OpenTelemetry tracer provider for use in Fiber middleware.
func Tracer() *sdktrace.TracerProvider {
	exporter, err := stdout.New()
	if err != nil {
		log.Fatal(err)
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		log.Fatal("cannot retrieve build info")
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("databuddy"),
				semconv.ServiceVersionKey.String(info.Main.Version),
			),
		),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}
