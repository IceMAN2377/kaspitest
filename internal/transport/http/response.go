package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"reflect"
)

func ResponseWithError(logger *slog.Logger, w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	resp := response{
		Success: false,
		Err:     msg,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Error("failed on encoding json error: " + err.Error())
	}
}

func Response(logger *slog.Logger, w http.ResponseWriter, msg any, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	if msg != nil {
		if reflect.ValueOf(msg).Kind() == reflect.Slice && reflect.ValueOf(msg).Len() == 0 {
			msg = []any{}
		}

		if err := json.NewEncoder(w).Encode(msg); err != nil {
			logger.Error("failed on encoding http msg: " + err.Error())
		}
	} else {
		if err := json.NewEncoder(w).Encode([]any{}); err != nil {
			logger.Error("failed on encoding http msg: " + err.Error())
		}
	}
}

type response struct {
	Success bool   `json:"success"`
	Err     string `json:"errors"`
}
