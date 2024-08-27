package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync"

	"interactive-presentation/src/models"
	"interactive-presentation/src/storage"
	"interactive-presentation/src/utilities"
)

func GetCurrentPoll(w http.ResponseWriter, r *http.Request) {
	presentationUUID, err := utilities.ParseUUIDFromRequest(r, "presentation_id")
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid presentation ID", http.StatusBadRequest)
		return
	}

	var presentations []models.PresentationDB
	err = storage.SelectFromTable("presentation", presentationUUID, &presentations)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("Error selecting from presentation table: %v", err), http.StatusInternalServerError)
		return
	}
	if len(presentations) == 0 {
		log.Println("No presentation found")
		http.Error(w, "No presentation found", http.StatusNotFound)
		return
	}

	var polls []models.PollDB
	err = storage.SelectFromTable("poll", presentationUUID, &polls)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("Error selecting from poll table: %v", err), http.StatusInternalServerError)
		return
	}

	var currentPoll models.Poll
	for _, poll := range polls {
		if presentations[0].CurrentPollIndex == poll.Index {
			var optionsDB []models.OptionDB
			err = storage.SelectFromTable("option", poll.PollID, &optionsDB)
			if err != nil {
				log.Println(err)
				http.Error(w, fmt.Sprintf("Error selecting from option table: %v", err), http.StatusInternalServerError)
				return
			}
			sort.Slice(optionsDB, func(i, j int) bool {
				return optionsDB[i].Index < optionsDB[j].Index
			})
			var options []models.Option
			for _, option := range optionsDB {
				options = append(options, models.Option{Key: option.Key, Value: option.Value})
			}
			currentPoll = models.Poll{PollID: poll.PollID, Question: poll.Question, Options: options}
			break
		}
	}

	_ = utilities.WriteJSONResponse(w, currentPoll)
}

func PutCurrentPoll(w http.ResponseWriter, r *http.Request) {
	presentationUUID, err := utilities.ParseUUIDFromRequest(r, "presentation_id")
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid presentation ID", http.StatusBadRequest)
		return
	}

	var presentations []models.PresentationDB
	err = storage.SelectFromTable("presentation", presentationUUID, &presentations)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("Error selecting from presentation table: %v", err), http.StatusInternalServerError)
		return
	}
	if len(presentations) == 0 {
		log.Println("No presentation found")
		http.Error(w, "No presentation found", http.StatusNotFound)
		return
	}

	nextPollIndex := presentations[0].CurrentPollIndex + 1

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = storage.UpdatePresentation(presentationUUID, nextPollIndex)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error updating presentation table: %v", err), http.StatusInternalServerError)
			return
		}
	}()
	wg.Wait()

	err = storage.SelectFromTable("presentation", presentationUUID, &presentations)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("Error selecting from presentation table: %v", err), http.StatusInternalServerError)
		return
	}

	var polls []models.PollDB
	err = storage.SelectFromTable("poll", presentationUUID, &polls)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("Error selecting from poll table: %v", err), http.StatusInternalServerError)
		return
	}

	var nextPoll models.Poll
	for _, poll := range polls {
		if poll.Index == nextPollIndex {
			var optionsDB []models.OptionDB
			err = storage.SelectFromTable("option", poll.PollID, &optionsDB)
			if err != nil {
				log.Println(err)
				http.Error(w, fmt.Sprintf("Error selecting from option table: %v", err), http.StatusInternalServerError)
				return
			}
			sort.Slice(optionsDB, func(i, j int) bool {
				return optionsDB[i].Index < optionsDB[j].Index
			})
			var options []models.Option
			for _, option := range optionsDB {
				options = append(options, models.Option{Key: option.Key, Value: option.Value})
			}
			nextPoll = models.Poll{PollID: poll.PollID, Question: poll.Question, Options: options}
			break
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = utilities.WriteJSONResponse(w, nextPoll)
}
