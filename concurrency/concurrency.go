package concurrency

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"todo/dataaccess"
)

var baseUrl = "http://localhost:3000"

type request struct {
	method   string
	url      string
	reqbody  string
	response int
}

func StartConcurrency() {

	if !isServerRunning() {
		slog.Warn("Server is not running - unable to continue")
		return
	}

	concurrentRequests(10)

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

func generateRequests(num int) []request {
	var requests []request
	ids := dataaccess.GetCurrentIDS()

	// Actions are GET/DELETE/PUT/POST
	// Set up POST's / create first

	for i := 0; i < num; i++ {
		var myReq request
		myReq.method = "POST"
		myReq.url = fmt.Sprintf("%s/%s/", baseUrl, "todo")
		myReq.reqbody = fmt.Sprintf(`{"description" : "%s%d", "status" : "NOT STARTED"}`, "this is body number ", i)
		requests = append(requests, myReq)
	}

	// generate PUT/GET/DELETES
	for i := 0; i < len(ids); i++ {
		// random number 1-3
		// Seed the random number generator
		rand.Seed(time.Now().UnixNano())

		// Generate a random number between 1 and 3
		randomNumber := rand.Intn(3) + 1
		var myReq request

		switch randomNumber {
		case 1:
			// get
			myReq.method = "GET"
			myReq.url = fmt.Sprintf("%s/%s/%d", baseUrl, "todo", ids[i])
			myReq.reqbody = ""

		case 2:
			// put
			myReq.method = "PUT"
			myReq.url = fmt.Sprintf("%s/%s", baseUrl, "todo")
			myReq.reqbody = fmt.Sprintf(`{"id" : "%d" , "description" : "%s", "status" : "NOT STARTED"}`, ids[i], "updated body")

		case 3:
			// delete
			myReq.method = "DELETE"
			myReq.url = fmt.Sprintf("%s/%s/%d", baseUrl, "todo", ids[i])
			myReq.reqbody = ""
		}
		// add in our request
		requests = append(requests, myReq)
	}
	// shuffle them up so we get an intermixed set of responses
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	rand := rand.New(source)
	rand.Shuffle(len(requests), func(i, j int) { requests[i], requests[j] = requests[j], requests[i] })
	return requests
}

// this will service our http requests
func worker(id int, reqs <-chan request, results chan<- string) {
	for r := range reqs {

		start := time.Now()
		slog.Info(fmt.Sprintf("Worker %d started job %v with method %s\n", id, r, r.method))
		switch r.method {
		case "GET", "DELETE":
			// get or delete strip the id, add a request path value and fire it in
			segments := strings.Split(r.url, "/")
			id := segments[len(segments)-1]

			req, err := http.NewRequest(r.method, r.url, nil)
			if err != nil {
				results <- fmt.Sprintf("Error creating %s request for url %s", r.method, r.url)
				return
			}

			// add the id
			req.SetPathValue("id", id)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				results <- fmt.Sprintf("Error sending request URL: %s, error: %v", req.URL, err)
				return
			}
			defer resp.Body.Close()

			// send back the body (none in case of delete) & status
			// Read and print the response
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				results <- fmt.Sprintf("Unable to extract body for request URL: %s, error: %v", req.URL, err)
				return
			}
			elapsed := time.Since(start)
			results <- fmt.Sprintf("Processed %s in %v, Response Length: %d, Response Status: %d", r.url, elapsed, len(body), resp.StatusCode)

		case "PUT", "POST":

			// Convert body to byte array
			jsonBody := []byte(r.reqbody)
			req, err := http.NewRequest(r.method, r.url, bytes.NewBuffer(jsonBody))
			if err != nil {
				results <- fmt.Sprintf("Error creating %s request for url %s", r.method, r.url)
				return
			}

			// Set headers
			req.Header.Set("Content-Type", "application/json")

			// fire in request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				results <- fmt.Sprintf("Error sending request URL : %s, error : %v", req.URL, err)
				return
			}
			defer resp.Body.Close()

			// send back the body (none in case of delete) & status
			// Read and print the response
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				results <- fmt.Sprintf("Unable to extract body for request URL: %s, error: %v", req.URL, err)
				return
			}
			elapsed := time.Since(start)
			results <- fmt.Sprintf("Processed %s in %v, Response Length: %d, Response Status: %d", r.url, elapsed, len(body), resp.StatusCode)
		}
	}
}

func concurrentRequests(numberToRun int) {

	// generate some requests
	requests := generateRequests(numberToRun)
	for _, req := range requests {
		slog.Info("request: %v", req)
	}

	//setup channels and size appropriately
	requestChan := make(chan request, len(requests))
	responseChan := make(chan string, len(requests))

	// set up workers
	numOfWorkers := 10
	for i := 1; i <= numOfWorkers; i++ {
		go worker(i, requestChan, responseChan)
	}

	// send our requests in
	for _, req := range requests {
		requestChan <- req
	}
	// close up
	close(requestChan)

	slog.Info(fmt.Sprintln("Showing responses...."))
	for i := 0; i < len(requests); i++ {
		slog.Info(fmt.Sprintf("Generated response %s", <-responseChan))
	}
}
