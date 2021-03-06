package internal

import (
	"encoding/json"
	"time"

	"github.com/gorilla/mux"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"

	"io"
	"log"
	"net/http"

	//"github.com/gorilla/mux"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/core"
)

func CreateAliasHandler(s *Server) http.Handler {
	svc := newCoreService(s)
	meter := otel.Meter("api-shortenit-v1")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dec := json.NewDecoder(r.Body)
		var req core.ShortenURLRequest
		err := dec.Decode(&req)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		recordRequestSize(meter, r, req, GetTag(s.cfg))

		res, err := svc.GetNewAlias(r.Context(), req)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		toResponse(res, w)
	})
}

func GetUrlHandler(s *Server) http.Handler {
	svc := newCoreService(s)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		url, err := svc.GetUrl(r.Context(), vars["alias"])

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		toResponse(&core.URLResponse{URL: url}, w)
	})
}

func GetTag(cfg *config.Config) *label.KeyValue {
	return &label.KeyValue{
		Key:   "api-shortenit-v1-commit-id",
		Value: label.StringValue(cfg.CommitSHA),
	}
}

func recordRequestSize(meter metric.Meter, r *http.Request, req core.ShortenURLRequest, label *label.KeyValue) {
	counter := metric.Must(meter).NewInt64Counter("api-shortenit-v1.create-alias.request-size.total")
	counter.Add(r.Context(), req.Size(), *label)
}

func newCoreService(s *Server) core.Service {
	svc := NewService(s.aliasSvc, s.userRepo, s.aliasRepo, s.cfg)
	return svc
}

func toResponse(data interface{}, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.Encode(data)
}

func InitSampleHandler(s *Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		s.userRepo.CreateUser(r.Context(), &core.User{
			ID:        primitive.NewObjectID(),
			Name:      "John D",
			Email:     "john.d@gmail.com",
			CreatedAt: time.Now(),
		})

		s.userRepo.CreateUser(r.Context(), &core.User{
			ID:        primitive.NewObjectID(),
			Name:      "Foo Bar",
			Email:     "foo.bar@gmail.com",
			CreatedAt: time.Now(),
		})

		s.userRepo.CreateUser(r.Context(), &core.User{
			ID:        primitive.NewObjectID(),
			Name:      "Tony Tr",
			Email:     "tony@gmail.com",
			CreatedAt: time.Now(),
		})

		w.Write([]byte("Done"))
	}
}
