package main

// imports
import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"todo/cli"
	"todo/web"
)

// All code split out into packages, this is now the main entrypoint
func main() {
	// Set up our logging, could write to a file here !
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true}))

	// Use this logger throughout the app
	slog.SetDefault(logger)
	// Create a traceId from google UUID & store it in a context

	// Flag values
	var startup string

	// use flags here to determine which bit of code we run
	flag.StringVar(&startup, "startup", "show", "Selected action")
	flag.Parse()

	// Startup can only be web/cli
	msg := fmt.Sprintf("Selected startup..%s", strings.ToLower(startup))
	slog.Info(msg)
	switch strings.ToLower(startup) {
	case "cli":

		cli.StartToDo()
	case "web":

		web.StartMux()
	default:
		msg = fmt.Sprintf("Unsupported startup type..%s passed to procedure", strings.ToLower(startup))
		slog.Info(msg)
	}
}
