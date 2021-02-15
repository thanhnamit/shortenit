package platform

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"
)

const (
	CtxApiKeyName = "api-key"
	CtxBasePath   = "base-path"
)

type ContextKey string

func NewGlobalHandler(handler http.Handler, operation string) func(w http.ResponseWriter, r *http.Request) {
	return toHandlerFunc(withCORS(withAPIKey(otelhttp.NewHandler(handler, operation))), operation)
}

func NewHealthHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	}
}

func toHandlerFunc(next http.Handler, operation string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		meter := otel.Meter("api-shortenit-v1")
		requestStartTime := time.Now()

		path := fmt.Sprintf("http://%s%s", r.Host, r.URL.Path)
		ctx := context.WithValue(r.Context(), ContextKey(CtxBasePath), path)
		next.ServeHTTP(w, r.WithContext(ctx))

		elapsedTime := time.Since(requestStartTime).Microseconds()
		recorder := metric.Must(meter).NewInt64ValueRecorder("api-shortenit-v1.latency")

		labels := []label.KeyValue{label.String("operation", operation), label.String("http-method", r.Method)}
		recorder.Record(r.Context(), elapsedTime, labels...)
	}
}

// withCORS set CORS headers for JS clients
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Location")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// withAPIKey checks and extract api key
func withAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Authorization")
		if key == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// record api key to context
		ctx := context.WithValue(r.Context(), ContextKey(CtxApiKeyName), key)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
