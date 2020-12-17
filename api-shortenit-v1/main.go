package main

import (
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/config"
)

func main() {
	cfg := config.NewAppConfig()
	server := internal.NewServer(cfg)
	server.Start()
}