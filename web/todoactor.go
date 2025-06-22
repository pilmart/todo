package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"todo/dataaccess"
	"todo/model"
	"todo/utils"
)

type Message struct {
	Context context.Context
	Action  string
	Payload string
	Reply   chan Response
}

// Response represents the actor's reply
type Response struct {
	Status  int
	Message string
}

// Actor definition represents a single actor handling http request/response pairs
type Actor struct {
	mailbox chan Message
}

// actor code starts - NewActor creates a new actor struct with a buffered mailbox channel
func NewActor(bufSize int) *Actor {
	return &Actor{mailbox: make(chan Message, bufSize)}
}

// Start sets up the processing loop over the mailbox to pick up the incoming
// messages and pass them to the handle messge function
func (a *Actor) Start() {
	go func() {
		for msg := range a.mailbox {
			fmt.Printf("Processing: Request action : %s ", msg.Action)
			traceID := msg.Context.Value("traceID")
			switch strings.ToLower(msg.Action) {

			case "get":
				// grab the incoming id, with error trap
				toDoID, err := strconv.Atoi(msg.Payload)
				if err != nil {
					// build an error string
					errMsg := fmt.Sprintf("Trace ID: %v Unable to convert ID: %s Error returned: %s", traceID, msg.Payload, err.Error())
					// log the error
					slog.Error(errMsg)
					msg.Reply <- Response{Status: http.StatusBadRequest, Message: errMsg}
					return
				}
				toDo, err := dataaccess.GetByID(msg.Context, toDoID)
				if err != nil {
					// build an error string
					errMsg := fmt.Sprintf("Trace ID: %v Unable to run GetByID with id : %s Error returned: %s", traceID, msg.Payload, err.Error())
					// log the error
					slog.Error(errMsg)
					msg.Reply <- Response{Status: http.StatusNotFound, Message: errMsg}
					return
				}
				// marshal the data, could also use new encoder here as well
				jsonData, err := json.Marshal(toDo)
				if err != nil {
					// build an error string
					// build an error string
					errMsg := fmt.Sprintf("Trace ID: %v Unable to marshal toDo item : %v Error returned: %s", traceID, toDo, err.Error())
					// log the error
					slog.Error(errMsg)
					msg.Reply <- Response{Status: http.StatusInternalServerError, Message: errMsg}
					return
				}
				msg.Reply <- Response{Status: http.StatusOK, Message: string(jsonData)}

			case "delete":
				// grab the incoming id, with error trap
				toDoID, err := strconv.Atoi(msg.Payload)
				if err != nil {
					// build an error string
					errMsg := fmt.Sprintf("Trace ID: %v Unable to convert ID: %s Error returned: %s", traceID, msg.Payload, err.Error())
					// log the error
					slog.Error(errMsg)
					msg.Reply <- Response{Status: http.StatusBadRequest, Message: errMsg}
					return
				}
				err = dataaccess.Delete(msg.Context, toDoID)
				if err != nil {
					// build an error string
					errMsg := fmt.Sprintf("Trace ID: %v Unable to delete record ID:: %s Error returned: %s", traceID, msg.Payload, err.Error())
					// log the error
					slog.Error(errMsg)
					msg.Reply <- Response{Status: http.StatusNotFound, Message: errMsg}
					return
				}
				msg.Reply <- Response{Status: http.StatusOK, Message: "Deleted ok"}

			case "put":
				// decode request body for now, probably needs much better error handling !!!
				var toDo model.ToDo

				// Convert string to JSON object
				err := json.Unmarshal([]byte(msg.Payload), &toDo)
				if err != nil {
					// build an error string
					errMsg := fmt.Sprintf("Trace ID: %v Error marshalling todo object %v, Error returned: %s", traceID, msg.Payload, err.Error())
					// log the error
					slog.Error(errMsg)
					msg.Reply <- Response{Status: http.StatusInternalServerError, Message: errMsg}
					return
				}

				// assume all incoming data is good so we don't have to perform checks against
				// field changes
				err = dataaccess.Update(msg.Context, toDo)
				if err != nil {
					errMsg := fmt.Sprintf("Trace ID: %v Error updating todo object %v, Error returned: %s", traceID, msg.Payload, err.Error())
					// log the error
					slog.Error(errMsg)
					msg.Reply <- Response{Status: http.StatusInternalServerError, Message: errMsg}
					return
				}

				msg.Reply <- Response{Status: http.StatusOK, Message: "updated ok"}

			case "post":
				// get our toDo Object
				var toDo model.ToDo

				// Convert string to JSON object
				err := json.Unmarshal([]byte(msg.Payload), &toDo)
				if err != nil {
					// build an error string
					errMsg := fmt.Sprintf("Trace ID: %v Error marshalling todo object %v, Error returned: %s", traceID, msg.Payload, err.Error())
					// log the error
					slog.Error(errMsg)
					msg.Reply <- Response{Status: http.StatusInternalServerError, Message: errMsg}
					return
				}
				// Check the description is good
				if len(toDo.Description) == 0 {
					// blank description, build an error string
					errMsg := fmt.Sprintf("Trace ID: %v Error : %v", traceID, errors.New("description cannot be blank"))
					// log the error
					slog.Error(errMsg)
					msg.Reply <- Response{Status: http.StatusBadRequest, Message: errMsg}
					return
				}

				// Check the status is good
				if !utils.ValidateStatus(toDo.Status) {
					// incorrect status
					errMsg := fmt.Sprintf("Trace ID: %v Error : incorrect status, must be one of %s ", traceID, utils.ShowPermittedStatuses())
					slog.Error(errMsg)
					msg.Reply <- Response{Status: http.StatusBadRequest, Message: errMsg}
					return
				}

				err = dataaccess.Create(msg.Context, toDo.Description, toDo.Status)
				if err != nil {
					// incorrect status
					errMsg := fmt.Sprintf("Trace ID: %v Error creating new toDo Item %v", traceID, toDo)
					slog.Error(errMsg)
					msg.Reply <- Response{Status: http.StatusInternalServerError, Message: errMsg}
					return
				}
				msg.Reply <- Response{Status: http.StatusOK, Message: "created ok"}
			}

		}
	}()
}

// Stop gracefully shuts down the Actor.
func (a *Actor) Stop() {
	close(a.mailbox)
}
