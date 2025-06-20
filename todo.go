package main

// imports
import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"todo/cli"
	"todo/constants"
	"todo/web"

	"github.com/google/uuid"
)

// All code split out into packages, this is now the main entrypoint
func main() {
	// Set up our logging, could write to a file here !
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: false}))

	// Use this logger throughout the app
	slog.SetDefault(logger)

	// Flag values
	var Id int
	var action string
	var description string
	var status string
	var service string

	/*
	 Initial flag setup to capture new todo information
	 id - only relevant for update/delete options
	 action - one of show/update/create/delete
	 description - default to empty string - relevant for create
	 status - default to not started - relevant for create
	*/

	traceID := uuid.NewString()
	ctx := context.WithValue(context.Background(), "traceID", traceID)

	flag.StringVar(&service, "service", "web", "Selected service")
	flag.StringVar(&action, "action", "show", "Selected action")
	flag.StringVar(&description, "description", " ", "Description of to do item")
	flag.StringVar(&status, "status", constants.StatusNotStarted, "Status of to do item")
	flag.IntVar(&Id, "Id", 0, "Mandatory for both update/delete actions")
	flag.Parse()

	// actions - Create / Show / Update / Delete
	msg := fmt.Sprintf("Selected service %s", strings.ToLower(service))
	slog.Info(msg)

	switch strings.ToLower(service) {
	case "web":
		// fire up web server
		web.StartMux()
	case "cli":
		// show all records no params needed
		cli.StartCLI(ctx, action, status, description, Id)
	case "concurrency":
		//concurrency.StartConcurrency()

	default:
		slog.Info(msg + " currently unsupported")
	}
}
