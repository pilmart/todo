package web

import (
	"fmt"
	"log/slog"
	"net/http"
)

// Handlers go here
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	// write data directly to response
	fmt.Fprintf(w, "got about page coming soon....")
	slog.Info("path about called")
}

func StartMux() {

	mux := http.NewServeMux()

	// register a simple about handler
	mux.HandleFunc("/about", aboutHandler)

	http.ListenAndServe("localhost:3000", mux)
}
