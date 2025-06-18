package concurrency

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

var baseUrl = "http://localhost:3000"

func StartConcurrency() {

	if !isServerRunning() {
		slog.Warn("Server is not running - unable to continue")
		return
	}

	// generate multiple concurrent creates(num)
	concurrentCreate(10)

	// Get A list of ids to use
	// ids := dataaccess.GetCurrentIDS()
	// slog.Info(fmt.Sprintf("Ids %v", ids))

	// generate multiple concurrent reads(num)

	// generate multiple updates(num)

	// generate multiple deletes(num)

}

// simple private function to make sure our server is up and listening
// before we hit it
func isServerRunning() bool {
	result := false
	// Check if the server is running by sending a request to the about page
	url := fmt.Sprintf("%s/%s", baseUrl, "about")
	fmt.Printf("about to run : %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Server is not running: %v\n", err)
	} else {
		fmt.Printf("Server is running! Response: %s\n", resp.Status)
		resp.Body.Close()
		result = true
	}
	return result
}

func concurrentCreate(numberToRun int) {

	bodies := generateRequests(numberToRun)
	processedBodies := processRequests(bodies)
	for str := range processedBodies {
		fmt.Println(str)
	}
}

// func concurrentRead(numberToRun int) {}

// func concurrentDelete(numberToRun int) {}

// func concurrentUpdate(numberToRun int) {}

func generateRequests(numberToRun int) <-chan string {

	out := make(chan string)

	// generate request bodies and add them to the channel
	go func() {

		for i := 0; i <= numberToRun; i++ {
			// stick the body on the channel
			body := fmt.Sprintf(`{"description" : "%s%d", "status" : "NOT STARTED"}`, "this is body number ", i)
			slog.Info("Generated", "body", body)
			out <- body
		}
		close(out)
	}()
	return out
}

func processRequests(in <-chan string) <-chan string {

	out := make(chan string)

	// generate request bodies and add them to the channel
	go func() {

		for body := range in {
			// fire in a create request
			url := fmt.Sprintf("%s/%s", baseUrl, "todo")
			// set the body
			reqBody := strings.NewReader(body)
			resp, _ := http.Post(url, "application/json", reqBody)
			out <- fmt.Sprintf("processed body : %s status code : %d", body, resp.StatusCode)
		}
		close(out)
	}()
	return out
}
