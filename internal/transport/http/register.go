package http

import (
	"log/slog"
	"net/http"

	"github.com/marioscordia/egov/internal/service"
)

func RegisterEndpoints(logger *slog.Logger, router *http.ServeMux, service service.Service) {
	handler := newHandler(service, logger)

	router.HandleFunc("GET /iin_check/{iin}", handler.CheckIIN)
	router.HandleFunc("GET /people/info/iin/{iin}", handler.GetByIIN)
	router.HandleFunc("GET /people/info/phone/{search}", handler.GetBySearch)
	router.HandleFunc("POST /people/info", handler.CreatePerson)
}
