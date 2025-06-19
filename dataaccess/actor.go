package dataaccess

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"todo/model"
)

// Actor represents a single actor handling requests
type Actor struct {
	mailbox chan Request
}

// Request represents a message sent to the actor
type Request struct {
	Context context.Context
	Method  string
	Path    string
	Body    string
	Reply   chan string
}

// NewActor creates a new actor with a buffered mailbox
func NewActor(bufferSize int) *Actor {
	return &Actor{
		mailbox: make(chan Request, bufferSize),
	}
}

// Start begins processing messages in the actor's mailbox
func (a *Actor) Start() {
	go func() {
		response := ""
		for req := range a.mailbox {
			// Process the request - need to check the incoming request method and fire up the
			// appropriate dataaccess method
			switch strings.ToLower(req.Method) {
			case "get":
				// retrieve
				// get or delete strip the id
				segments := strings.Split(req.Path, "/")
				intId, _ := strconv.Atoi(segments[len(segments)-1])
				toDo, err := GetByID(req.Context, intId)
				if err != nil {
					// build an error string
					errMsg := fmt.Sprintf("Error: Unable to run GetByID with id : %d Error returned: %s TraceID: %v", intId, err.Error(), req.Context.Value("traceID"))
					req.Reply <- errMsg
					// log the error
					slog.Error(errMsg)
					return
				}

				// Convert the struct to a JSON string
				jsonData, err := json.Marshal(toDo)
				if err != nil {
					// build an error string
					errMsg := fmt.Sprintf("Error: Unable to marshal data for id : %d Error returned: %s TraceID: %v", intId, err.Error(), req.Context.Value("traceID"))
					req.Reply <- errMsg
					// log the error
					slog.Error(errMsg)
					return
				}
				// hand back data
				//req.Reply <- fmt.Sprintf("TraceId: %v request method: %s request path:", req.Context.Value("traceID"), req.Method, req.Path)
				response = string(jsonData)
			case "post":
				// create
				// Decode the JSON body into the struct
				var toDo model.ToDo
				reader := strings.NewReader(req.Body)
				err := json.NewDecoder(reader).Decode(&toDo)
				if err != nil {
					// build an error string
					errMsg := fmt.Sprintf("Error: Unable to decode data for creation using req body %v, Error returned: %s TraceID: %v", req.Body, err.Error(), req.Context.Value("traceID"))
					req.Reply <- errMsg
					// log the error
					slog.Error(errMsg)
					return
				}

				err = Create(req.Context, toDo.Description, toDo.Status)
				if err != nil {
					// build an error string
					errMsg := fmt.Sprintf("Error: Creating new toDo Item %s TraceId: %v", err.Error(), req.Context.Value("traceID"))
					req.Reply <- errMsg
					// log the error
					slog.Error(errMsg)
					return
				}
				response = fmt.Sprintf("TraceId: %v request method: %s request path: %s", req.Context.Value("traceID"), req.Method, req.Path)
			case "put":
				// update
				// Decode the JSON body into the struct
				var toDo model.ToDo
				reader := strings.NewReader(req.Body)
				err := json.NewDecoder(reader).Decode(&toDo)
				if err != nil {
					// build an error string
					errMsg := fmt.Sprintf("Error: Unable to decode data for creation using req body %v, Error returned: %s TraceID: %v", req.Body, err.Error(), req.Context.Value("traceID"))
					req.Reply <- errMsg
					// log the error
					slog.Error(errMsg)
					return
				}

				err = Update(req.Context, toDo)
				if err != nil {
					// build an error string
					errMsg := fmt.Sprintf("Error: Creating new toDo Item %s TraceId: %v", err.Error(), req.Context.Value("traceID"))
					req.Reply <- errMsg
					// log the error
					slog.Error(errMsg)
					return
				}
				response = fmt.Sprintf("TraceId: %v request method: %s request path: %s", req.Context.Value("traceID"), req.Method, req.Path)
			case "delete":
				// delete
				segments := strings.Split(req.Path, "/")
				intId, _ := strconv.Atoi(segments[len(segments)-1])
				err := Delete(req.Context, intId)
				if err != nil {
					// build an error string
					errMsg := fmt.Sprintf("Error: Unable to delete record ID: %d Error returned: %s TraceId: %v", intId, err.Error(), req.Context.Value("traceID"))
					req.Reply <- errMsg
					// log the error
					slog.Error(errMsg)
					return

				}
				response = fmt.Sprintf("TraceId: %v request method: %s request path: %s", req.Context.Value("traceID"), req.Method, req.Path)
			}
			req.Reply <- response
		}
	}()
}

// HandleRequest sends a request to the actor mailbox
func (a *Actor) HandleRequest(req Request) {
	a.mailbox <- req
}
