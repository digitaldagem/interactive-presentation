package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"interactive-presentation/src/models"
	"interactive-presentation/src/storage"
	"interactive-presentation/src/utilities"
)

func PostPollVote(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := utilities.ReadRequestBody(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var vote models.Vote
	if err = json.Unmarshal(bodyBytes, &vote); err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err = storage.InsertIntoDatabase("vote", vote); err != nil {
		log.Println(err)
		http.Error(w, "Error inserting into vote database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetPollVotes(w http.ResponseWriter, r *http.Request) {
	pollUUID, err := utilities.ParseUUIDFromRequest(r, "poll_id")
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid poll ID", http.StatusBadRequest)
		return
	}

	var votes []models.Vote
	err = storage.SelectFromTable("vote", pollUUID, &votes)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("Error selecting votes: %v", err), http.StatusInternalServerError)
		return
	}

	_ = utilities.WriteJSONResponse(w, &votes)
}
