package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/robinmuhia/helm-charts/pkg/helm-charts/application/common"
	"github.com/robinmuhia/helm-charts/pkg/helm-charts/application/helpers"
	"github.com/robinmuhia/helm-charts/pkg/helm-charts/presentation"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func setupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}

		shutdownFuncs = nil

		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)

	// Set up trace provider.
	tracerProvider, err := newTraceProvider()
	if err != nil {
		handleErr(err)
		return
	}

	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)

	otel.SetTracerProvider(tracerProvider)

	return
}

func newTraceProvider() (*trace.TracerProvider, error) {
	environment, err := helpers.GetEnvVar(common.Environment.String())
	if err != nil {
		return nil, err
	}

	serviceName := fmt.Sprintf("helm-charts-%v", environment)

	jaegerEndpoint, err := helpers.GetEnvVar(common.JaegerCollectorEndpoint.String())
	if err != nil {
		return nil, err
	}

	traceExporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(jaegerEndpoint),
	)
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(
			traceExporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultScheduleDelay*time.Millisecond),
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
		),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
			),
		),
	)

	_ = traceProvider.Tracer("helm-charts")

	return traceProvider, nil
}

// StartApplication is used to start the application server
func StartApplication(ctx context.Context) error {
	otelShutdown, err := setupOTelSDK(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err = errors.Join(err, otelShutdown(ctx))
	}()

	portEnv, err := helpers.GetEnvVar(common.Port.String())
	if err != nil {
		return err
	}

	port, err := strconv.Atoi(portEnv)
	if err != nil {
		return err
	}

	return presentation.StartServer(ctx, port)
}

func main() {
	ctx := context.Background()

	err := StartApplication(ctx)
	if err != nil {
		panic(fmt.Errorf("unable to start application: %w", err))
	}
}
