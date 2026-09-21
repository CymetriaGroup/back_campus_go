package server

import (
	"net/http"

	"hexagonal-go-backend/internal/platform/config"
)

func New(handler http.Handler, config config.AppConfig) *http.Server {
	return &http.Server{Addr: ":" + config.Port, Handler: handler, ReadHeaderTimeout: 5_000_000_000, ReadTimeout: config.ReadTimeout, WriteTimeout: config.WriteTimeout, IdleTimeout: config.IdleTimeout}
}
