package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"todo/dataaccess"
	"todo/model"
	"todo/utils"

	"github.com/google/uuid"
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

// Handlers go here similar to controllers in spring mvc
// Returns a static about page no need for mutex here....probably
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	//attempt to use new web/static/about.htm
	templateFile := "web/static/about.htm"

	// Trace ID should have been bolted on via anonymous wrapper
	ctx := r.Context()
	traceID := ctx.Value("traceID")

	slog.Info("path about called", "traceID", traceID)
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		http.Error(w, "Unable to load template file", http.StatusInternalServerError)
		slog.Info("Unable to load ", "template", templateFile)
		return
	}
	// now to execute the template
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Unable to execute template file", http.StatusInternalServerError)
		slog.Info("Unable to execute ", "template", templateFile)
		return
	}
	slog.Info("path about completes")
}

// Returns a 'slightly' more dynamic todo listing page
func toDoListHandler(w http.ResponseWriter, r *http.Request) {

	//attempt to use new web/static/about.htm
	templateFile := "web/templates/todolist.htm"

	// Trace ID should have been bolted on via anonymous wrapper
	ctx := r.Context()
	traceID := ctx.Value("traceID")

	slog.Info("todos path called ", "traceID", traceID)
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		http.Error(w, "Unable to load template file", http.StatusInternalServerError)
		slog.Info("Unable to load ", "template", templateFile)
		return
	}

	// Go and grab the current set of todos...
	toDos := dataaccess.GetAllRecords(ctx)

	// So far so good, pass some data similar to modelandview add object
	data := struct {
		Title   string
		Heading string
		TraceID string
		ToDos   []model.ToDo
	}{
		Title:   "Todos Page",
		Heading: "Current To Do Listing",
		TraceID: traceID.(string),
		ToDos:   toDos,
	}

	// now to execute the template
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Unable to execute template file", http.StatusInternalServerError)
		slog.Info("Unable to execute ", "template", templateFile)
		return
	}
	slog.Info("path todos completes")
}

// Simple middleware code to 'bolt-on' a traceID to the handlers
func withTraceID(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a traceId from google UUID & store it in a context
		ctx := context.WithValue(r.Context(), "traceID", uuid.NewString())
		r = r.WithContext(ctx)
		f(w, r)
	}
}

func StartMux() {

	actor := NewActor(5)
	actor.Start()

	mux := http.NewServeMux()

	slog.Info("Server mux started, available at http://localhost:3000")

	mux.HandleFunc("GET /about", withTraceID(aboutHandler))
	mux.HandleFunc("GET /todolist", withTraceID(toDoListHandler))

	// use anonymous functions so we can access actor & mailbox
	mux.HandleFunc("GET /todo/{id}", withTraceID(func(w http.ResponseWriter, r *http.Request) {
		// Trace ID should have been bolted on via withTraceID
		ctx := r.Context()
		traceID := ctx.Value("traceID")

		// set up reply channel
		reply := make(chan Response)
		defer close(reply)

		slog.Info(fmt.Sprintf("TraceID: %v Revised get Handler starts...", traceID))

		// create our message and send it directly to the mailbox
		actor.mailbox <- Message{Context: ctx, Action: r.Method, Payload: r.PathValue("id"), Reply: reply}

		// get our response
		response := <-reply

		// tell the user we're handing back json and set status 200 OK
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.Status)
		// probably easier ways to do this
		w.Write([]byte(response.Message))
		slog.Info(fmt.Sprintf("TraceID: %v Revised get Handler completes", traceID))

	}))

	mux.HandleFunc("DELETE /todo/{id}", withTraceID(func(w http.ResponseWriter, r *http.Request) {
		// Trace ID should have been bolted on via withTraceID
		ctx := r.Context()
		traceID := ctx.Value("traceID")

		// set up reply channel
		reply := make(chan Response)
		defer close(reply)

		slog.Info(fmt.Sprintf("TraceID: %v Revised delete Handler starts...", traceID))

		// create our message and send it directly to the mailbox
		actor.mailbox <- Message{Context: ctx, Action: r.Method, Payload: r.PathValue("id"), Reply: reply}

		// get our response
		response := <-reply

		w.WriteHeader(response.Status)
		// probably easier ways to do this
		w.Write([]byte(response.Message))
		slog.Info(fmt.Sprintf("TraceID: %v Revised delete Handler completes", traceID))

	}))

	mux.HandleFunc("PUT /todo", withTraceID(func(w http.ResponseWriter, r *http.Request) {
		// Trace ID should have been bolted on via withTraceID
		ctx := r.Context()
		traceID := ctx.Value("traceID")

		// set up reply channel
		reply := make(chan Response)
		defer close(reply)

		slog.Info(fmt.Sprintf("TraceID: %v Revised update Handler starts...", traceID))

		// convert request body to a string
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close() // Ensure the body is closed after reading

		// Convert the body to a string
		bodyString := string(bodyBytes)

		// create our message and send it directly to the mailbox
		actor.mailbox <- Message{Context: ctx, Action: r.Method, Payload: bodyString, Reply: reply}

		// get our response
		response := <-reply
		w.WriteHeader(response.Status)
		// probably easier ways to do this
		w.Write([]byte(response.Message))
		slog.Info(fmt.Sprintf("TraceID: %v Revised update Handler completes", traceID))

	}))

	mux.HandleFunc("POST/todo", withTraceID(func(w http.ResponseWriter, r *http.Request) {
		// Trace ID should have been bolted on via withTraceID
		ctx := r.Context()
		traceID := ctx.Value("traceID")

		// set up reply channel
		reply := make(chan Response)
		defer close(reply)

		slog.Info(fmt.Sprintf("TraceID: %v Revised create Handler starts...", traceID))

		// convert request body to a string
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close() // Ensure the body is closed after reading

		// Convert the body to a string
		bodyString := string(bodyBytes)

		// create our message and send it directly to the mailbox
		actor.mailbox <- Message{Context: ctx, Action: r.Method, Payload: bodyString, Reply: reply}

		// get our response
		response := <-reply
		w.WriteHeader(response.Status)
		// probably easier ways to do this
		w.Write([]byte(response.Message))
		slog.Info(fmt.Sprintf("TraceID: %v Revised create Handler completes", traceID))

	}))

	http.ListenAndServe("localhost:3000", mux)
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
			fmt.Printf("Processing: Request action :%s", msg.Action)
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
