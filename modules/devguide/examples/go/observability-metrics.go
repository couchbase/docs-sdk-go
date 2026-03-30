package main

import (
	"log"
	"time"

	gocbopentelemetry "github.com/couchbase/gocb-opentelemetry"
	"github.com/couchbase/gocb/v2"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"context"
)

func loggingMeter() {
	// tag::logging-meter[]
	meter := gocb.NewLoggingMeter(&gocb.LoggingMeterOptions{
		EmitInterval: 30 * time.Second,
	})
	// end::logging-meter[]

	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: "Administrator",
			Password: "password",
		},
		Meter: meter,
	}

	cluster, err := gocb.Connect("localhost", opts)
	if err != nil {
		log.Fatal(err)
	}

	cluster.Close(nil)
}

func openTelemetryMeter() {
	// tag::otel-meter[]
	ctx := context.Background()

	// Setup an exporter.
	// This exporter exports metrics on the OTLP protocol over GRPC to localhost:4317.
	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpointURL("http://localhost:4317"),
		otlpmetricgrpc.WithCompressor("gzip"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create the OpenTelemetry SDK's MeterProvider.
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			// An OpenTelemetry service name generally reflects the name of your microservice,
			// e.g. "shopping-cart-service".
			semconv.ServiceNameKey.String("YOUR_SERVICE_NAME_HERE"),
		)),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(1*time.Second))),
	)
	defer meterProvider.Shutdown(ctx)

	// Provide the MeterProvider to the Couchbase OpenTelemetry meter wrapper.
	meter := gocbopentelemetry.NewOpenTelemetryMeter(meterProvider)

	// Provide the OpenTelemetry meter as part of the Cluster configuration.
	cluster, err := gocb.Connect("localhost", gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: "Administrator",
			Password: "password",
		},
		Meter: meter,
	})
	if err != nil {
		log.Fatal(err)
	}
	// end::otel-meter[]

	cluster.Close(nil)
}
