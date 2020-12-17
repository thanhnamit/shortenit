package internal

import (
	"encoding/json"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/core"
	"io"
	"log"
	"net/http"
)

func CreateAliasHandler(s *Server) http.Handler {
	svc := NewService(s.aliasSvc, s.userRepo, s.cfg)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dec := json.NewDecoder(r.Body)
		var req core.ShortenURLRequest
		err := dec.Decode(&req)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		res, err := svc.NewAlias(r.Context(), req)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		toResponse(res, w)
	})
}

func GetUrlHandler(s *Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}


func toResponse(data interface{}, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.Encode(data)
}