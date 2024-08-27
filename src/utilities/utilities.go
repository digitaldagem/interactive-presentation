package utilities

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func ParseUUIDFromRequest(r *http.Request, paramName string) (uuid.UUID, error) {
	paramValue := chi.URLParam(r, paramName)
	return uuid.Parse(paramValue)
}

func WriteJSONResponse(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return err
	}
	return nil
}

func ReadRequestBody(r *http.Request) ([]byte, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %v", err)
	}
	var bodyCloseError error
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			err = bodyCloseError
		}
	}(r.Body)
	return bodyBytes, bodyCloseError
}
