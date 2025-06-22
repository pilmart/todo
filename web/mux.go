package web

import (
	"context"

	"fmt"
	"io"
	"log/slog"
	"net/http"
	"text/template"
	"todo/dataaccess"
	"todo/model"

	"github.com/google/uuid"
)

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

// Service 'GET' Requests
func getHandler(actor *Actor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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

	}
}

// Service 'PUT' / update Requests
func putHandler(actor *Actor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

// Service 'POST' / create Requests
func postHandler(actor *Actor) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

// Service 'DELETE'Requests
func deleteHandler(actor *Actor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
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

	// set up handlers
	mux.HandleFunc("GET /about", withTraceID(aboutHandler))
	mux.HandleFunc("GET /todolist", withTraceID(toDoListHandler))
	// rest calls
	mux.HandleFunc("GET /todo/{id}", withTraceID(getHandler(actor)))
	mux.HandleFunc("PUT /todo", withTraceID(putHandler(actor)))
	mux.HandleFunc("DELETE /todo/{id}", withTraceID(deleteHandler(actor)))
	mux.HandleFunc("POST /todo", withTraceID(postHandler(actor)))
	http.ListenAndServe("localhost:3000", mux)
}
