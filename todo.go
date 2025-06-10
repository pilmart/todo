package main

// imports
import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"strings"
	"todo/constants"
	"todo/dataaccess"
	"todo/datatypes"
	"todo/utils"

	"github.com/google/uuid"
)

func main() {

	// Flag values
	var Id int
	var action string
	var description string
	var status string

	// Set up our logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true}))

	slog.SetDefault(logger)
	// Create a traceId from google UUID & store it in a context
	traceID := uuid.NewString()
	ctx := context.WithValue(context.Background(), "traceID", traceID)
	ctx = context.WithValue(ctx, "random", "abcde")

	// Log with contextual information
	logger.With(
		"traceID", ctx.Value("traceID"),
		"random", ctx.Value("random"),
	).Info("test of ctx")

	/*
	 Initial flag setup to capture new todo information
	 id - only relevant for update/delete options
	 action - one of show/update/create/delete
	 description - default to empty string - relevant for create
	 status - default to not started - relevant for create
	*/

	flag.StringVar(&action, "action", "show", "Selected action")
	flag.StringVar(&description, "description", " ", "Description of to do item")
	flag.StringVar(&status, "status", constants.StatusNotStarted, "Status of to do item")
	flag.IntVar(&Id, "Id", 0, "Mandatory for both update/delete actions")
	flag.Parse()

	// actions - Create / Show / Update / Delete
	log.Printf("Selected action..%s\n", strings.ToLower(action))
	switch strings.ToLower(action) {
	case "show":
		// show all records no params needed
		dataaccess.ShowAllRecords()
	case "create":
		if !utils.ValidateStatus(status) {
			utils.ShowPermittedStatuses()
			log.Fatalf("Status of %s is not permitted", status)
			return
		}
		// All ok create new todo item
		dataaccess.Create(description, status)
	case "update":
		var updatedToDo datatypes.ToDo
		updatedToDo.Id = Id
		updatedToDo.Description = description
		updatedToDo.Status = status
		if !utils.ValidateStatus(status) {
			utils.ShowPermittedStatuses()
			log.Fatalf("Status of %s is not permitted", status)
			return
		}
		// all ok update may continue
		dataaccess.Update(updatedToDo)
	case "delete":
		// delete function
		dataaccess.Delete(Id)
	default:
		log.Printf("Unsupported action..%s passed to procedure\n", strings.ToLower(action))
	}
}
