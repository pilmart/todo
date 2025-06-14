package main

// imports
import (
	"log/slog"
	"os"
	"todo/cli"
)

// All code split out into packages, this is now the main entrypoint
func main() {
	// Set up our logging, could write to a file here !
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true}))

	// Use this logger throughout the app
	slog.SetDefault(logger)

	//web.StartMux()
	cli.StartToDo()

}
