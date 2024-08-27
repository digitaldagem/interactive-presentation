# Interactive Presentation

A backend service that serves poll data and poll vote data for an interactive presentation app.
The services stores polls and poll votes in a postgres database and presentations are fetched from `https://infra.devskills.app/api/interactive-presentation/v4` every time they are needed.

### Endpoints

* test endpoint to see if service is up and running
  * `GET /ping`
* proxy endpoint that makes a call for a presentation to be created
  * `POST /presentations`
* endpoint to fetch data for the current poll
  * `GET /presentations/{presentation_id}/polls/current`
* endpoint to update the current poll index for a presentation and get the next poll
  * `PUT /presentations/{presentation_id}/polls/current`
* endpoint to record a poll vote
  * `POST /presentations/{presentation_id}/polls/current/votes`
* endpoint to fetch all votes for a given poll
  * `GET /presentations/{presentation_id}/polls/{poll_id}/votes`

### Running the service locally in docker
Run `make up-local`  
### Testing the service while it is running locally in docker
Run `make test`
### Stopping the service and removing its images in docker
Run `make down`  