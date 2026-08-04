package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // 1 MegaByte
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(data)
}

func ResponseJSONError(w http.ResponseWriter, status int, message string) error {
	type envelop struct {
		Error string `json:"error"`
	}
	return WriteJSON(w, status, &envelop{Error: message})
}

func ResponseJSON(w http.ResponseWriter, status int, data any) error {
	type responseJSON struct {
		Data any `json:"data"`
	}
	return WriteJSON(w, status, &responseJSON{Data: data})
}
