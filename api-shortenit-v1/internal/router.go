package internal

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/platform"
)

type Router struct {
	Handler http.Handler
}

func NewRouter(s *Server) *Router {
	r := mux.NewRouter()
	r.HandleFunc("/init-sample-data", InitSampleHandler(s))
	r.HandleFunc("/shortenit", platform.NewGlobalHandler(CreateAliasHandler(s), "POST /shortenit"))
	r.HandleFunc("/shortenit/{alias}", platform.NewGlobalHandler(GetUrlHandler(s), "GET /shortenit/{alias}"))
	r.HandleFunc("/health", platform.NewHealthHandler())
	return &Router{Handler: r}
}
