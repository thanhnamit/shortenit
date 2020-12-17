package internal

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/core"
	"io"
	"log"
	"net/http"
)

func CreateAliasHandler(s *Server) http.Handler {
	svc := NewService(s.aliasSvc, s.userRepo, s.aliasRepo, s.cfg)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dec := json.NewDecoder(r.Body)
		var req core.ShortenURLRequest
		err := dec.Decode(&req)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		res, err := svc.GetNewAlias(r.Context(), req)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		toResponse(res, w)
	})
}

func GetUrlHandler(s *Server) http.Handler {
	svc := NewService(s.aliasSvc, s.userRepo, s.aliasRepo, s.cfg)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		url, err := svc.GetUrl(r.Context(), vars["alias"])

		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		toResponse(&core.URLResponse{URL: url}, w)
	})
}


func toResponse(data interface{}, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.Encode(data)
}