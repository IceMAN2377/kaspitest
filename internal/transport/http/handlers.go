package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	errs "github.com/marioscordia/egov/internal/errors"
	"github.com/marioscordia/egov/internal/models"
	"github.com/marioscordia/egov/internal/service"
)

func newHandler(service service.Service, logger *slog.Logger) *handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

type handler struct {
	service service.Service
	logger  *slog.Logger
}

func (h *handler) CheckIIN(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	iin := r.PathValue("iin")

	info, err := h.service.CheckIIN(ctx, iin)
	if err != nil {
		ResponseWithError(h.logger, w, err.Error(), http.StatusBadRequest)
		return
	}

	Response(h.logger, w, info, http.StatusOK)
}

func (h *handler) GetByIIN(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	iin := r.PathValue("iin")

	user, err := h.service.GetByIIN(ctx, iin)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			ResponseWithError(h.logger, w, err.Error(), http.StatusNotFound)
			return
		}
		ResponseWithError(h.logger, w, err.Error(), http.StatusInternalServerError)
		return
	}

	Response(h.logger, w, user, http.StatusOK)
}

func (h *handler) GetBySearch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	search := r.PathValue("search")

	users, err := h.service.GetBySearch(ctx, search)
	if err != nil {
		if !errors.Is(err, errs.ErrNotFound) {
			ResponseWithError(h.logger, w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	Response(h.logger, w, users, http.StatusOK)
}

func (h *handler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		ResponseWithError(h.logger, w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.IIN == "" || user.Name == "" || user.Phone == "" {
		ResponseWithError(h.logger, w, "IIN, name and phone required", http.StatusBadRequest)
		return
	}

	if err := h.service.CreatePerson(ctx, &user); err != nil {
		if errors.Is(err, errs.ErrAlreadyExists) || errors.Is(err, errs.ErrIncorrectData) {
			ResponseWithError(h.logger, w, err.Error(), http.StatusBadRequest)
			return
		}

		ResponseWithError(h.logger, w, err.Error(), http.StatusInternalServerError)
		return
	}

	Response(h.logger, w, nil, http.StatusOK)
}
