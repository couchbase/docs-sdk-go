package main

import (
	"context"
	"log"
	"time"

	gocbopentelemetry "github.com/couchbase/gocb-opentelemetry"
	"github.com/couchbase/gocb/v2"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func main() {}

func thresholdLoggingTracer() {
	// tag::threshold-logging[]
	tracer := gocb.NewThresholdLoggingTracer(&gocb.ThresholdLoggingOptions{
		Interval:    1 * time.Minute,
		SampleSize:  10,
		KVThreshold: 2 * time.Second,
	})
	// end::threshold-logging[]

	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: "Administrator",
			Password: "password",
		},
		Tracer: tracer,
	}

	cluster, err := gocb.Connect("localhost", opts)
	if err != nil {
		log.Fatal(err)
	}

	cluster.Close(nil)
}

func openTelemetryTracer() {
	// tag::otel-tracing[]
	ctx := context.Background()

	// Setup an exporter.
	// This exporter exports traces on the OTLP protocol over GRPC to localhost:4317.
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpointURL("http://localhost:4317"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create the OpenTelemetry SDK's TracerProvider.
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			// An OpenTelemetry service name generally reflects the name of your microservice,
			// e.g. "shopping-cart-service".
			semconv.ServiceNameKey.String("YOUR_SERVICE_NAME_HERE"),
		)),
		// The BatchSpanProcessor will efficiently batch traces and periodically export them.
		sdktrace.WithBatcher(exporter),
		// Export every trace: this may be too heavy for production.
		// An alternative is `sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.01))`
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	defer tracerProvider.Shutdown(ctx)

	// Provide the TracerProvider to the Couchbase OpenTelemetry tracer wrapper.
	tracer := gocbopentelemetry.NewOpenTelemetryRequestTracer(tracerProvider)

	// Provide the OpenTelemetry tracer as part of the Cluster configuration.
	cluster, err := gocb.Connect("localhost", gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: "Administrator",
			Password: "password",
		},
		Tracer: tracer,
	})
	if err != nil {
		log.Fatal(err)
	}
	// end::otel-tracing[]

	bucket := cluster.Bucket("travel-sample")
	col := bucket.DefaultCollection()

	// tag::parent-span[]
	// Create an OpenTelemetry span to use as a parent.
	otelCtx, parentOtelSpan := tracerProvider.Tracer("my-app").Start(ctx, "my-operation")
	defer parentOtelSpan.End()

	// Wrap it in a Couchbase RequestSpan and pass it as a parent to the SDK operation.
	parentSpan := gocbopentelemetry.NewOpenTelemetryRequestSpan(otelCtx, parentOtelSpan)

	result, err := col.Get("my-doc", &gocb.GetOptions{
		ParentSpan: parentSpan,
	})
	// end::parent-span[]
	_, _ = result, err

	cluster.Close(nil)
}
