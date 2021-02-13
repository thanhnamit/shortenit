package tracing

import (
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	otellabel "go.opentelemetry.io/otel/label"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	jaegerUrl = "http://localhost:14268/api/traces"
)

// InitTracer ...
func InitTracer(serviceName string) func() {

	url := os.Getenv("TRACER_COLLECTOR")
	if url == "" {
		url = jaegerUrl
	}

	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint(url),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: serviceName,
			Tags: []otellabel.KeyValue{
				otellabel.String("exporter", "jaeger"),
			},
		}),
		jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
	)
	if err != nil {
		log.Fatal(err)
	}

	// required for trace context to be propagated
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return flush
}
