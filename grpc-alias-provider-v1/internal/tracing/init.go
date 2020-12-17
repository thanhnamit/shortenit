package tracing

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/propagators"
	"log"

	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	otellabel "go.opentelemetry.io/otel/label"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	jaegerUrl = "http://localhost:14268/api/traces"
)

// InitTracer ...
func InitTracer(serviceName string) func() {
	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint(jaegerUrl),
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
	global.SetTextMapPropagator(otel.NewCompositeTextMapPropagator(propagators.TraceContext{}, propagators.Baggage{}))
	return flush
}
