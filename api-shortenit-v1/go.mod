module github.com/thanhnamit/shortenit/api-shortenit-v1

go 1.15

require (
	go.mongodb.org/mongo-driver v1.4.3
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.13.0
	go.opentelemetry.io/otel v0.13.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.13.0
	go.opentelemetry.io/otel/sdk v0.13.0
)
