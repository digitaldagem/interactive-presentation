package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"

	"interactive-presentation/src/config"
	"interactive-presentation/src/handlers"
)

const (
	Presentation = "presentation"
	Poll         = "poll"
	Option       = "option"
	Vote         = "vote"
)

func pingHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Service is up and running"))
	if err != nil {
		log.Println("error writing data for pingHandler: ", err)
	}
}

func main() {
	configuration, err := config.New()
	if err != nil {
		log.Fatal("error creating and initializing a new configuration object: ", err)
	}

	db, err := newDB(configuration)
	if err != nil {
		log.Fatal("error connecting to database: ", err)
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			log.Println("error closing database connection: ", err)
		}
	}(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/ping", pingHandler)

	r.Post("/presentations", handlers.CreatePresentation)

	r.Get("/presentations/{presentation_id}/polls/current", handlers.GetCurrentPoll)
	r.Put("/presentations/{presentation_id}/polls/current", handlers.PutCurrentPoll)

	r.Post("/presentations/{presentation_id}/polls/current/votes", handlers.PostPollVote)
	r.Get("/presentations/{presentation_id}/polls/{poll_id}/votes", handlers.GetPollVotes)

	log.Println("Starting server on :8080...")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("error listening on the TCP network address and calling Serve with handler", err)
	}
}

func newDB(config *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error opening a database connection: %v", err))
	}

	statement := "CREATE TABLE IF NOT EXISTS " + Presentation + " (presentation_id uuid PRIMARY KEY, current_poll_index integer);"
	_, err = db.Exec(statement)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error creating %s table: %v", Presentation, err))
	}

	statement = "CREATE TABLE IF NOT EXISTS " + Poll + " (poll_id uuid PRIMARY KEY, question VARCHAR(255), presentation_id uuid, index integer);"
	_, err = db.Exec(statement)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error creating %s table: %v", Poll, err))
	}

	statement = "CREATE TABLE IF NOT EXISTS " + Option + " (key VARCHAR(255), value VARCHAR(255), poll_id uuid, index integer);"
	_, err = db.Exec(statement)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error creating %s table: %v", Option, err))
	}

	statement = "CREATE TABLE IF NOT EXISTS " + Vote + " (key VARCHAR(255), client_id VARCHAR(255), poll_id uuid);"
	_, err = db.Exec(statement)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error creating %s table: %v", Vote, err))
	}

	log.Println("Connected to database")

	return db, nil
}
