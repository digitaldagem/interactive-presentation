package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"interactive-presentation/src/models"
)

const apiUrl = "http://localhost:8080"

func TestBackend(t *testing.T) {

	t.Run("calling ping should return 200", func(t *testing.T) {
		resp, err := http.Get(apiUrl + "/ping")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "The application should respond with 200 on ping request")
	})

	t.Run("creating a presentation with invalid data should result with 400 status code", func(t *testing.T) {
		resp, err := http.Post(apiUrl+"/presentations", "application/json", bytes.NewBuffer([]byte("{}")))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Response status should be 400 when creating a template with no polls")
	})

	t.Run("reading a presentation using an unknown presentation id should result with 404 status code", func(t *testing.T) {
		req, err := http.NewRequest("GET", apiUrl+"/presentations/425ba663-ecaf-4902-84c2-2f7d1aa3d1d7", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode, "A 404 status code should be returned in case of presentation of given id not found")
	})

	t.Run("creating a presentation, showing a poll, voting and reading votes", func(t *testing.T) {
		var presentationID string
		var pollID uuid.UUID

		// Create a presentation
		body := models.Presentation{
			CurrentPollIndex: 0,
			Polls: []models.Poll{
				{
					Question: "What's your favorite pet?",
					Options: []models.Option{
						{Key: "A", Value: "Dog"},
						{Key: "B", Value: "Cat"},
						{Key: "C", Value: "Crocodile"},
					},
				},
				{
					Question: "Which of the countries would you like to visit the most?",
					Options: []models.Option{
						{Key: "A", Value: "Argentina"},
						{Key: "B", Value: "Austria"},
						{Key: "C", Value: "Australia"},
					},
				},
			},
		}

		bodyBytes, err := json.Marshal(body)
		assert.NoError(t, err)

		resp, err := http.Post(apiUrl+"/presentations", "application/json", bytes.NewBuffer(bodyBytes))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Response status is 201 when creating a presentation")

		var result map[string]string
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)
		presentationID = result["presentation_id"]
		assert.NotEmpty(t, presentationID, "Presentation ID is present when creating a presentation")

		// Get the current poll
		resp, err = http.Get(apiUrl + "/presentations/" + presentationID + "/polls/current")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "getting current poll should return 200 status code")

		var poll models.Poll
		err = json.NewDecoder(resp.Body).Decode(&poll)
		assert.NoError(t, err)
		pollID = poll.PollID
		assert.NotEmpty(t, pollID, "poll_id should be returned")
		assert.Equal(t, "What's your favorite pet?", poll.Question)

		// Present the next poll
		req, err := http.NewRequest("PUT", apiUrl+"/presentations/"+presentationID+"/polls/current", nil)
		assert.NoError(t, err)
		client := &http.Client{}
		resp, err = client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "switching to the next poll must return 200 http code")

		var poll2 models.Poll
		err = json.NewDecoder(resp.Body).Decode(&poll2)
		assert.NoError(t, err)
		pollID = poll2.PollID
		assert.NotEmpty(t, pollID, "poll_id should be returned")
		assert.Equal(t, "Which of the countries would you like to visit the most?", poll2.Question)

		// Get the votes for the current poll
		resp, err = http.Get(apiUrl + "/presentations/" + presentationID + "/polls/" + pollID.String() + "/votes")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "reading presentation votes should return 200")

		var votes []models.Vote
		err = json.NewDecoder(resp.Body).Decode(&votes)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(votes), "votes list should be empty")

		// Vote for the current poll
		clientID := uuid.NewString()
		key, err := generateRandomCapitalLetter()
		assert.NoError(t, err)
		voteBody := models.Vote{
			Key:      key,
			ClientID: clientID,
			PollID:   pollID,
		}

		voteBodyBytes, err := json.Marshal(voteBody)
		assert.NoError(t, err)

		req, err = http.NewRequest("POST", apiUrl+"/presentations/"+presentationID+"/polls/current/votes", bytes.NewBuffer(voteBodyBytes))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		resp, err = client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode, "voting should return 204 status code")

		// Verify the vote was recorded
		resp, err = http.Get(apiUrl + "/presentations/" + presentationID + "/polls/" + pollID.String() + "/votes")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "reading votes should return 200 status code")

		err = json.NewDecoder(resp.Body).Decode(&votes)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(votes), "votes list should contain one vote")
		expectedVote := models.Vote{
			ClientID: clientID,
			Key:      key,
		}
		assert.Equal(t, expectedVote.Key, votes[0].Key, "the recorded key should match the expected key")
		assert.Equal(t, expectedVote.ClientID, votes[0].ClientID, "the recorded client_id should match the expected client_id")
	})
}

func generateRandomCapitalLetter() (string, error) {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
	if err != nil {
		return "", err
	}
	return string(letters[n.Int64()]), nil
}
