package platform

import (
	"context"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
)

const (
	CtxApiKeyName = "api-key"
)

type ContextKey string

func NewGlobalHandler(handler http.Handler, operation string) func(w http.ResponseWriter, r *http.Request) {
	return toHandlerFunc(withCORS(withAPIKey(otelhttp.NewHandler(handler, operation))))
}

func NewHealthHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	}
}

func toHandlerFunc(next http.Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
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

