package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"

	"interactive-presentation/src/models"
	"interactive-presentation/src/storage"
	"interactive-presentation/src/utilities"
)

const (
	baseURL = "https://infra.devskills.app/api/interactive-presentation/v4"
)

func CreatePresentation(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := utilities.ReadRequestBody(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url := baseURL + "/presentations"
	resp, err := http.Post(url, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create presentation", http.StatusInternalServerError)
		return
	}
	var bodyCloseError error
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			err = bodyCloseError
		}
	}(resp.Body)
	if err != nil {
		log.Println("error closing response body", err)
	}

	if resp.StatusCode != http.StatusCreated {
		http.Error(w, "Failed to create presentation", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	var result map[string]interface{}
	if err = json.Unmarshal(body, &result); err != nil {
		http.Error(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}

	var presentation models.Presentation
	if err = json.Unmarshal(bodyBytes, &presentation); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var presentationDB models.PresentationDB
	presentationID, ok := result["presentation_id"].(string)
	presentationUUID, err := uuid.Parse(presentationID)
	if ok {
		presentationDB.PresentationID = presentationUUID
		presentationDB.CurrentPollIndex = 0
	}

	if err = storage.InsertIntoDatabase("presentation", presentationDB); err != nil {
		http.Error(w, fmt.Sprintf("Error inserting into presentation database: %v", err), http.StatusInternalServerError)
		return
	}

	for i, poll := range presentation.Polls {
		pollID := uuid.New()
		pollDB := models.PollDB{PollID: pollID, Question: poll.Question, PresentationID: presentationUUID, Index: i}
		if err = storage.InsertIntoDatabase("poll", pollDB); err != nil {
			http.Error(w, fmt.Sprintf("Error inserting into poll database: %v", err), http.StatusInternalServerError)
			return
		}

		for j, option := range poll.Options {
			optionDB := models.OptionDB{Key: option.Key, Value: option.Value, PollID: pollID, Index: j}
			if err = storage.InsertIntoDatabase("option", optionDB); err != nil {
				http.Error(w, fmt.Sprintf("Error inserting into option database: %v", err), http.StatusInternalServerError)
				return
			}
		}
	}

	w.WriteHeader(resp.StatusCode)
	_, err = w.Write(body)
	if err != nil {
		log.Println(err)
	}
}
