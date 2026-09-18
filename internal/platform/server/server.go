package server

import (
	"hexagonal-go-backend/internal/config"
	"net/http"
)

func New(handler http.Handler, c config.AppConfig) *http.Server {
	return &http.Server{Addr: ":" + c.Port, Handler: handler, ReadHeaderTimeout: 5_000_000_000, ReadTimeout: c.ReadTimeout, WriteTimeout: c.WriteTimeout, IdleTimeout: c.IdleTimeout}
}
