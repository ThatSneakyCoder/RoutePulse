package main

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // maximum one 1mB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	// stop unknown fields
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func writeJSONError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}

	return writeJSON(w, status, &envelope{Error: message})
}

func (app *application) jsonResponse(w http.ResponseWriter, status int, data ...any) error {
	if len(data) == 0 {
		w.WriteHeader(status)
		return nil
	}

	type envelope struct {
		Data any `json:"data"`
	}

	return writeJSON(w, status, &envelope{data})
}

func (app *application) httpResponse(w http.ResponseWriter, status int, message string) error {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)

	if message != "" {
		_, err := w.Write([]byte(message))
		if err != nil {
			return err
		}
	}

	return nil
}
