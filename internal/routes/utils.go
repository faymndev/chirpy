package routes

import (
	"encoding/json"
	"net/http"
)

func SendJSON(w http.ResponseWriter, statusCode int, payload map[string]any) error {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Write(data)
	return nil
}
