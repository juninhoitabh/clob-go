package shared

import (
	"encoding/json"
	"net/http"
)

func Mul(a, b int64) int64 { return a * b }

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
