package routes

import (
	"encoding/json"
	"net/http"
)

func SendJSON[T any](w http.ResponseWriter, statusCode int, payload T) error {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Write(data)
	return nil
}
