package web

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"todo/dataaccess"
	"todo/model"

	"github.com/google/uuid"
)

var mu sync.Mutex

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

	// synchronisation here
	mu.Lock()
	defer mu.Unlock()

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

// Get a single record
func getHandler(w http.ResponseWriter, r *http.Request) {

	// synchronisation here
	mu.Lock()
	defer mu.Unlock()

	// Trace ID should have been bolted on via anonymous wrapper
	ctx := r.Context()
	traceID := ctx.Value("traceID")

	slog.Info("Rest Call to getHandler", "TraceID", traceID)
	// grab the incoming id, with error trap
	toDoID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		// build an error string
		errMsg := fmt.Sprintf("Unable to convert ID: %s Error returned: %s TraceID: %v", r.PathValue("id"), err.Error(), traceID)
		// log the error
		slog.Error(errMsg)
		// return http error + status
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	slog.Info("Incoming id", "ToDO ID :", toDoID)

	// Go And grab our toDo Item
	toDo, err := dataaccess.GetByID(ctx, toDoID)
	if err != nil {
		// build an error string
		errMsg := fmt.Sprintf("Unable to execute GetByID with id : %s Error returned: %s TraceID: %v", r.PathValue("id"), err.Error(), traceID)
		// log the error
		slog.Error(errMsg)
		// return http error + status
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
	// marshal the data, could also use new encoder here as well
	jsonData, err := json.Marshal(toDo)
	if err != nil {
		// build an error string
		errMsg := fmt.Sprintf("Error marshalling todo object, Error returned: %s TraceID: %v", err.Error(), traceID)
		// log the error
		slog.Error(errMsg)
		// return http error + status
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	// tell the user we're handing back json and set status 200 OK
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
	slog.Info("Rest Call to getHandler completes - json data sent")

}

// delete a single record
func deleteHandler(w http.ResponseWriter, r *http.Request) {

	// synchronisation here
	mu.Lock()
	defer mu.Unlock()

	// Trace ID should have been bolted on via anonymous wrapper
	ctx := r.Context()
	traceID := ctx.Value("traceID")

	slog.Info("Rest Call to deleteHandler", "traceID", traceID)
	// grab the incoming id, with error trap
	toDoID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		// build an error string
		errMsg := fmt.Sprintf("Unable to convert ID: %s Error returned: %s TraceId: %v", r.PathValue("id"), err.Error(), traceID)
		// log the error
		slog.Error(errMsg)
		// return http error + status
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	slog.Info("Incoming id", "ToDO ID :", toDoID)

	// Go and delete our toDo Item
	err = dataaccess.Delete(ctx, toDoID)
	if err != nil {
		// build an error string
		errMsg := fmt.Sprintf("Unable to delete record ID: %s Error returned: %s TraceId: %v", r.PathValue("id"), err.Error(), traceID)
		// log the error
		slog.Error(errMsg)
		// return http error + status
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	// all good hand back status nocontent 204
	w.WriteHeader(http.StatusNoContent)
	slog.Info("Rest Call to deleteHandler completes")

}

// Update a todo item, requires a model.ToDo payload using 'PUT' method
func updateHandler(w http.ResponseWriter, r *http.Request) {

	// synchronisation here
	mu.Lock()
	defer mu.Unlock()

	// Trace ID should have been bolted on via anonymous wrapper
	ctx := r.Context()
	traceID := ctx.Value("traceID")

	slog.Info("Rest Call to updateHandler", "traceID", traceID)

	// decode request body for now, probably needs much better error handling !!!
	var toDo model.ToDo
	err := json.NewDecoder(r.Body).Decode(&toDo)
	if err != nil {
		// build an error string
		errMsg := fmt.Sprintf("Error marshalling todo object, Error returned: %s TraceId: %v", err.Error(), traceID)
		// log the error
		slog.Error(errMsg)
		// return http error + status
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	// get original record to compare
	var originalToDo model.ToDo
	originalToDo, err = dataaccess.GetByID(ctx, toDo.Id)
	if err != nil {
		// build an error string
		errMsg := fmt.Sprintf("Error locating (original) todo with id %d, Error returned: %s TraceId: %v", toDo.Id, err.Error(), traceID)
		// log the error
		slog.Error(errMsg)
		// return http error + status
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
	// no description present use original
	if len(toDo.Description) == 0 {
		//use original
		toDo.Description = originalToDo.Description
	}

	// no status present use original
	if len(toDo.Status) == 0 {
		//use original
		toDo.Status = originalToDo.Status
	}

	// call the update
	err = dataaccess.Update(ctx, toDo)
	if err != nil {
		// build an error string
		errMsg := fmt.Sprintf("Error updating todo object %v, Error returned: %s TraceId: %v", toDo, err.Error(), traceID)
		// log the error
		slog.Error(errMsg)
		// return http error + status
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
	// all good hand back ok
	w.WriteHeader(http.StatusOK)
	slog.Info("updateHandler done")

}

func createHandler(w http.ResponseWriter, r *http.Request) {
	// synchronisation here
	mu.Lock()
	defer mu.Unlock()

	// Trace ID should have been bolted on via anonymous wrapper
	ctx := r.Context()
	traceID := ctx.Value("traceID")

	slog.Info("Rest Call to createHandler", "traceID", traceID)

	// decode request body for now, probably needs much better error handling !!!
	var toDo model.ToDo
	err := json.NewDecoder(r.Body).Decode(&toDo)
	if err != nil {
		// build an error string
		errMsg := fmt.Sprintf("Error marshalling todo object, Error returned: %s TraceId: %v", err.Error(), traceID)
		// log the error
		slog.Error(errMsg)
		// return http error + status
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	// capture description
	description := toDo.Description

	// capture the status
	status := toDo.Status

	// all ok pass description & status to create function
	// Additional validation occurs in create
	dataaccess.Create(ctx, description, status)
	if err != nil {
		// build an error string
		errMsg := fmt.Sprintf("Error creating new toDo Item %s TraceId: %v", err.Error(), traceID)
		// log the error
		slog.Error(errMsg)
		// return http error + status
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	// all good hand back ok
	w.WriteHeader(http.StatusOK)
	slog.Info("createHandler done")

}

func StartMux() {

	mux := http.NewServeMux()

	fmt.Println("Server mux started, available at http://localhost:3000")
	fmt.Println("about page can be found at http://localhost:3000/about")
	fmt.Println("todolist page can be found at http://localhost:3000/todolist")
	fmt.Println("GET /todo/{id}, returns single record")
	fmt.Println("DELETE /todo/{id}, deletes a single record identified by {id}")
	fmt.Println("PUT /todo, with request body, updates a record")
	fmt.Println("POST /todo, with request body, creates a new record")

	// Handler registration
	// register a simple about handler, added trace ID - implemented
	mux.HandleFunc("GET /about", func(w http.ResponseWriter, r *http.Request) {
		// Create a traceId from google UUID & store it in a context
		ctx := context.WithValue(r.Context(), "traceID", uuid.NewString())
		r = r.WithContext(ctx)
		aboutHandler(w, r)
	})

	// register a simple todo list handler use this to show a list of current todos - implemented
	mux.HandleFunc("GET /todolist", func(w http.ResponseWriter, r *http.Request) {
		// Create a traceId from google UUID & store it in a context
		ctx := context.WithValue(r.Context(), "traceID", uuid.NewString())
		r = r.WithContext(ctx)
		toDoListHandler(w, r)
	})

	// register get by id handler - implemented
	mux.HandleFunc("GET /todo/{id}", func(w http.ResponseWriter, r *http.Request) {
		// Create a traceId from google UUID & store it in a context
		ctx := context.WithValue(r.Context(), "traceID", uuid.NewString())
		r = r.WithContext(ctx)
		getHandler(w, r)
	})

	// register delete handler - implemented
	mux.HandleFunc("DELETE /todo/{id}", func(w http.ResponseWriter, r *http.Request) {
		// Create a traceId from google UUID & store it in a context
		ctx := context.WithValue(r.Context(), "traceID", uuid.NewString())
		r = r.WithContext(ctx)
		deleteHandler(w, r)
	})

	// register update handler - implemented
	mux.HandleFunc("PUT /todo", func(w http.ResponseWriter, r *http.Request) {
		// Create a traceId from google UUID & store it in a context
		ctx := context.WithValue(r.Context(), "traceID", uuid.NewString())
		r = r.WithContext(ctx)
		updateHandler(w, r)
	})

	// register create handler
	mux.HandleFunc("POST /todo", func(w http.ResponseWriter, r *http.Request) {
		// Create a traceId from google UUID & store it in a context
		ctx := context.WithValue(r.Context(), "traceID", uuid.NewString())
		r = r.WithContext(ctx)
		createHandler(w, r)
	})

	http.ListenAndServe("localhost:3000", mux)
}
