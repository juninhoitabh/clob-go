package shared

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Errors struct {
	Status  int      `json:"status" example:"400"`
	Message string   `json:"message" example:"Invalid parameter"`
	Details []string `json:"details,omitempty" example:"The 'name' field is required"`
}

type ErrorResponse struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

func Mul(a, b int64) int64 { return a * b }

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, err error, status int) {
	errResp := ErrorResponse{
		Status:  status,
		Message: err.Error(),
	}

	WriteJSON(w, status, errResp)
}

func HandleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		WriteError(w, err, http.StatusNotFound)
	case errors.Is(err, ErrAlreadyExists):
		WriteError(w, err, http.StatusConflict)
	case errors.Is(err, ErrInvalidParam):
		WriteError(w, err, http.StatusBadRequest)
	default:
		WriteError(w, err, http.StatusInternalServerError)
	}
}

func BadRequestError(w http.ResponseWriter, message string, details ...string) {
	errResp := ErrorResponse{
		Status:  http.StatusBadRequest,
		Message: message,
		Details: details,
	}

	WriteJSON(w, http.StatusBadRequest, errResp)
}
