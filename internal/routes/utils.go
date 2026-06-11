package routes

import (
	"encoding/json"
	"errors"
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

func DecodeBody[T any](r *http.Request) (T, error) {
	var input T

	if r.Body == nil {
		return input, errors.New("Request body is empty")
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	return input, err
}

