package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	tracing "github.com/thanhnamit/shortenit/api-shortenit-v1/infra/tracing"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/model"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/label"
)

var userRepo *user.Repository

func main() {
	ctx := context.Background()
	userRepo = user.NewRepository(ctx)
	defer userRepo.Close(ctx)

	flush := tracing.InitTracer("api-shortenit-v1")
	defer flush()

	// instrumented
	http.Handle("/admin/users", otelhttp.NewHandler(http.HandlerFunc(handleGetUsers), "/admin/users"))
	http.Handle("/shortenit/", otelhttp.NewHandler(http.HandlerFunc(handleShortenIt), "/shortenit"))

	// not instrumented
	http.HandleFunc("/admin/init", handleInitUsers)

	// start server
	log.Println("Listening on 8085...")
	log.Fatal(http.ListenAndServe(":8085", nil))
}

func handleShortenIt(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleRedirectURL(w, r)
	case http.MethodPost:
		handleCreateURL(w, r)
	case http.MethodDelete:
		handleDeleteURL(w, r)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func handleRedirectURL(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/shortenit/")
	w.Write([]byte(key))
}

func handleCreateURL(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var req model.ShortenURLRequest
	err := dec.Decode(&req)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//
	toJSON(req, w)
}

func handleDeleteURL(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/shortenit/")
	// delete by key
	w.Write([]byte(key))
}

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	// create span with global named context
	ctx := r.Context()
	tr := global.Tracer("api-shortenit-v1")
	ctx, span := tr.Start(ctx, "main.get-users")
	defer span.End()

	users, err := userRepo.GetAllUsers(ctx)
	if err != nil {
		span.AddEvent(ctx, "db.error", label.String("message", err.Error()))
		log.Fatal(err)
	}

	re, err := json.Marshal(users)
	if err != nil {
		span.AddEvent(ctx, "convert.error", label.String("message", err.Error()))
		panic(err)
	}

	// add attributes (tags in opentracing) - for query
	span.SetAttributes(label.String("database-type", "mongodb"))

	// add event (logs in opentracing) - good to record specific, context-based information
	span.AddEvent(ctx, "response.size", label.Int("size", len(users)))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(re)
}

func handleInitUsers(w http.ResponseWriter, r *http.Request) {
	userRepo.CreateUser(r.Context(), &user.User{
		ID:        primitive.NewObjectID(),
		Name:      "Nam1",
		Email:     "thanhnam.it@gmail.com",
		CreatedAt: time.Now(),
	})

	userRepo.CreateUser(r.Context(), &user.User{
		ID:        primitive.NewObjectID(),
		Name:      "Nam2",
		Email:     "thanhnam2.it@gmail.com",
		CreatedAt: time.Now(),
	})

	userRepo.CreateUser(r.Context(), &user.User{
		ID:        primitive.NewObjectID(),
		Name:      "Nam3",
		Email:     "thanhnam3.it@gmail.com",
		CreatedAt: time.Now(),
	})
}

func toJSON(data interface{}, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.Encode(data)
}
