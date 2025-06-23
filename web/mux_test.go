package web

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

// test endpoint "GET /todo/{id}"
func TestGetHandler(t *testing.T) {

	baseUrl := "/todo"

	tests := []struct {
		testName     string
		id           int
		expectedCode int
		expectedBody string
	}{
		{"happy path record present", 34, http.StatusOK, `{"id":34,"description":"new record via bruno on sunday":"NOT STARTED"}`},
		{"non-existent record", 445, http.StatusNotFound, `{}`},
		{"id zero test", 0, http.StatusNotFound, `{}`},
		{"id minus one test", -1, http.StatusNotFound, `{}`},
	}

	for _, test := range tests {

		actor := NewActor(1)
		actor.Start()
		urlUnderTest := fmt.Sprintf("%s/%d", baseUrl, test.id)
		t.Logf("Url under test : %s, test name GET - %s", urlUnderTest, test.testName)
		req := httptest.NewRequest(http.MethodGet, urlUnderTest, nil)
		// add on the path variable... we hope
		req.SetPathValue("id", strconv.Itoa(test.id))
		rec := httptest.NewRecorder()
		withTraceID(getHandler(actor)).ServeHTTP(rec, req)
		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != test.expectedCode {
			t.Errorf("For URL %s, expected status %d, got %d", urlUnderTest, test.expectedCode, res.StatusCode)
		}

		body := rec.Body.String()
		t.Logf("returned body : %s", body)
		actor.Stop()
	}
}

// parallel get requests
func TestGetHandlerParallel(t *testing.T) {

	baseUrl := "/todo"

	tests := []struct {
		testName     string
		id           int
		expectedCode int
		expectedBody string
	}{

		{"happy path record present", 25, http.StatusOK, `{"id":25,"description":"this is body number 2":"NOT STARTED"}`},
		{"happy path record present", 26, http.StatusOK, `{"id":26,"description":"this is body number 3":"NOT STARTED"}`},
		{"happy path record present", 27, http.StatusOK, `{"id":27,"description":"this is body number 4":"NOT STARTED"}`},
		{"happy path record present", 28, http.StatusOK, `{"id":28,"description":"this is body number 5":"NOT STARTED"}`},
		{"happy path record present", 29, http.StatusOK, `{"id":29,"description":"this is body number 6":"NOT STARTED"}`},
		{"happy path record present", 18, http.StatusOK, `{"id":18,"description":"Updated via test":"NOT STARTED"}`},
		{"happy path record present", 23, http.StatusOK, `{"id":23,"description":"this is body number 0":"NOT STARTED"}`},
		{"happy path record present", 20, http.StatusOK, `{"id":20,"description":"Updated via test":"NOT STARTED"}`},
		{"happy path record present", 21, http.StatusOK, `{"id":21,"description":"Updated via update test":"NOT STARTED"}`},
		{"happy path record present", 22, http.StatusOK, `{"id":22,"description":"Updated via test":"NOT STARTED"}`},
		{"happy path record present", 34, http.StatusOK, `{"id":34,"description":"new record via bruno on sunday":"NOT STARTED"}`},
		{"non-existent record", 445, http.StatusNotFound, `{}`},
		{"id zero test", 0, http.StatusNotFound, `{}`},
		{"id minus one test", -1, http.StatusNotFound, `{}`},
	}

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			actor := NewActor(10)
			actor.Start()
			t.Parallel()
			urlUnderTest := fmt.Sprintf("%s/%d", baseUrl, test.id)
			t.Logf("Url under test : %s, test name GET - %s", urlUnderTest, test.testName)
			req := httptest.NewRequest(http.MethodGet, urlUnderTest, nil)
			req.SetPathValue("id", strconv.Itoa(test.id))
			rec := httptest.NewRecorder()
			withTraceID(getHandler(actor)).ServeHTTP(rec, req)
			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != test.expectedCode {
				t.Errorf("For URL %s, expected status %d, got %d", urlUnderTest, test.expectedCode, res.StatusCode)
			}
			actor.Stop()
		})

	}

}

// Test Endpoint "POST /todo"
func TestPostHandler(t *testing.T) {

	baseUrl := "/todo"
	tests := []struct {
		testName     string
		expectedCode int
		expectedBody string
	}{
		{"happy path", http.StatusOK, `{"description":"Updated via test","status": "NOT STARTED"}`},
		{"missing description", http.StatusBadRequest, `{"description":"","status": "NOT STARTED"}`},
		{"missing status", http.StatusBadRequest, `{"description":"Updated via test","status": ""}`},
		{"Incorrect status", http.StatusBadRequest, `{"description":"Updated via test","status": "Incorrect"}`},
		{"empty body", http.StatusBadRequest, `{}`},
	}

	for _, test := range tests {
		actor := NewActor(1)
		actor.Start()
		urlUnderTest := baseUrl
		t.Logf("Url under test : %s, test name POST - %s", urlUnderTest, test.testName)
		reqBody := strings.NewReader(test.expectedBody)
		req := httptest.NewRequest(http.MethodPost, urlUnderTest, reqBody)
		rec := httptest.NewRecorder()
		withTraceID(postHandler(actor)).ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != test.expectedCode {
			t.Errorf("For URL %s, expected status %d, got %d", urlUnderTest, test.expectedCode, res.StatusCode)
		}

		body := rec.Body.String()
		t.Logf("returned body : %s", body)
		actor.Stop()
	}
}

// Parallel Post requests
func TestPostHandlerParallel(t *testing.T) {

	baseUrl := "/todo"

	tests := []struct {
		testName     string
		expectedCode int
		expectedBody string
	}{

		{"happy path", http.StatusOK, `{"description":"created via test - 1","status": "NOT STARTED"}`},
		{"missing description", http.StatusBadRequest, `{"description":"","status": "NOT STARTED"}`},
		{"missing status", http.StatusBadRequest, `{"description":"Updated via test","status": ""}`},
		{"Incorrect status", http.StatusBadRequest, `{"description":"Updated via test","status": "Incorrect"}`},
		{"empty body", http.StatusBadRequest, `{}`},
		{"happy path", http.StatusOK, `{"description":"created via test - 2","status": "NOT STARTED"}`},
		{"missing description", http.StatusBadRequest, `{"description":"","status": "NOT STARTED"}`},
		{"missing status", http.StatusBadRequest, `{"description":"Updated via test","status": ""}`},
		{"Incorrect status", http.StatusBadRequest, `{"description":"Updated via test","status": "Incorrect"}`},
		{"empty body", http.StatusBadRequest, `{}`},
		{"happy path", http.StatusOK, `{"description":"created via test - 3","status": "NOT STARTED"}`},
		{"missing description", http.StatusBadRequest, `{"description":"","status": "NOT STARTED"}`},
		{"missing status", http.StatusBadRequest, `{"description":"Updated via test","status": ""}`},
		{"Incorrect status", http.StatusBadRequest, `{"description":"Updated via test","status": "Incorrect"}`},
		{"empty body", http.StatusBadRequest, `{}`},
		{"happy path", http.StatusOK, `{"description":"created via test - 4","status": "NOT STARTED"}`},
		{"missing description", http.StatusBadRequest, `{"description":"","status": "NOT STARTED"}`},
		{"missing status", http.StatusBadRequest, `{"description":"Updated via test","status": ""}`},
		{"Incorrect status", http.StatusBadRequest, `{"description":"Updated via test","status": "Incorrect"}`},
		{"empty body", http.StatusBadRequest, `{}`},
	}

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			actor := NewActor(10)
			actor.Start()
			t.Parallel()
			urlUnderTest := baseUrl
			t.Logf("Url under test : %s, test name POST - %s", urlUnderTest, test.testName)
			reqBody := strings.NewReader(test.expectedBody)
			req := httptest.NewRequest(http.MethodPost, urlUnderTest, reqBody)
			rec := httptest.NewRecorder()
			withTraceID(postHandler(actor)).ServeHTTP(rec, req)
			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != test.expectedCode {
				t.Errorf("For URL %s, expected status %d, got %d", urlUnderTest, test.expectedCode, res.StatusCode)
			}
			actor.Stop()
		})

	}

}

// test endpoint PUT /todo - need a current record
func TestPutHandler(t *testing.T) {

	baseUrl := "/todo"
	tests := []struct {
		testName     string
		expectedCode int
		expectedBody string
	}{
		{"happy path", http.StatusOK, `{"id":25, "description":"Updated via update test XXXABCDE","status": "NOT STARTED"}`},
		{"missing description", http.StatusOK, `{"id":21,"description":"","status": "NOT STARTED"}`},
		{"Incorrect status", http.StatusInternalServerError, `{"id":21,"description":"Updated via test","status": "Incorrect"}`},
		{"empty body", http.StatusInternalServerError, `{}`},
		{"non-existent id", http.StatusInternalServerError, `{"id":99999,"description":"Updated via update test","status": "NOT STARTED"}`},
	}

	for _, test := range tests {
		actor := NewActor(1)
		actor.Start()
		urlUnderTest := baseUrl
		t.Logf("Url under test : %s, test name %s", urlUnderTest, test.testName)
		reqBody := strings.NewReader(test.expectedBody)
		req := httptest.NewRequest(http.MethodPut, urlUnderTest, reqBody)
		rec := httptest.NewRecorder()
		withTraceID(putHandler(actor)).ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != test.expectedCode {
			t.Errorf("Test name %s, For URL %s, expected status %d, got %d", test.testName, urlUnderTest, test.expectedCode, res.StatusCode)
		}

		body := rec.Body.String()
		t.Logf("returned body : %s", body)
		actor.Stop()
	}
}

func TestPutHandlerParallel(t *testing.T) {

	baseUrl := "/todo"
	tests := []struct {
		testName     string
		expectedCode int
		expectedBody string
	}{
		{"happy path", http.StatusOK, `{"id":22, "description":"Updated via update test XXX 22","status": "NOT STARTED"}`},
		{"missing description", http.StatusInternalServerError, `{"id":29,"description":"","status": "NOT STARTED"}`},
		{"Incorrect status", http.StatusInternalServerError, `{"id":30,"description":"Updated via test","status": "Incorrect"}`},
		{"empty body", http.StatusInternalServerError, `{}`},
		{"non-existent id", http.StatusInternalServerError, `{"id":99999,"description":"Updated via update test","status": "NOT STARTED"}`},
		{"happy path", http.StatusOK, `{"id":23, "description":"Updated via update test XXX 23","status": "NOT STARTED"}`},
		{"missing description", http.StatusInternalServerError, `{"id":32,"description":"","status": "NOT STARTED"}`},
		{"Incorrect status", http.StatusInternalServerError, `{"id":33,"description":"Updated via test","status": "Incorrect"}`},
		{"empty body", http.StatusInternalServerError, `{}`},
		{"non-existent id", http.StatusInternalServerError, `{"id":99999,"description":"Updated via update test","status": "NOT STARTED"}`},
		{"happy path", http.StatusOK, `{"id":25, "description":"Updated via update test XXX 25","status": "NOT STARTED"}`},
		{"missing description", http.StatusInternalServerError, `{"id":21,"description":"","status": "NOT STARTED"}`},
		{"Incorrect status", http.StatusInternalServerError, `{"id":21,"description":"Updated via test","status": "Incorrect"}`},
		{"empty body", http.StatusInternalServerError, `{}`},
		{"non-existent id", http.StatusInternalServerError, `{"id":99999,"description":"Updated via update test","status": "NOT STARTED"}`},
		{"happy path", http.StatusOK, `{"id":26, "description":"Updated via update test XXX 26","status": "NOT STARTED"}`},
		{"missing description", http.StatusInternalServerError, `{"id":21,"description":"","status": "NOT STARTED"}`},
		{"Incorrect status", http.StatusInternalServerError, `{"id":21,"description":"Updated via test","status": "Incorrect"}`},
		{"empty body", http.StatusInternalServerError, `{}`},
		{"non-existent id", http.StatusInternalServerError, `{"id":99999,"description":"Updated via update test","status": "NOT STARTED"}`},
	}

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			actor := NewActor(10)
			actor.Start()
			t.Parallel()
			urlUnderTest := baseUrl
			t.Logf("Url under test : %s, test name %s", urlUnderTest, test.testName)
			reqBody := strings.NewReader(test.expectedBody)
			req := httptest.NewRequest(http.MethodPut, urlUnderTest, reqBody)
			rec := httptest.NewRecorder()
			withTraceID(putHandler(actor)).ServeHTTP(rec, req)
			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != test.expectedCode {
				t.Errorf("Test name %s, For URL %s, expected status %d, got %d", test.testName, urlUnderTest, test.expectedCode, res.StatusCode)
			}
			actor.Stop()
		})
	}
}

// test endpoint DELETE /todo/{id}
func TestDeleteHandler(t *testing.T) {
	t.Skip("Skipping this test for now")
	baseUrl := "/todo"
	tests := []struct {
		testName     string
		id           int
		expectedCode int
		expectedBody string
	}{
		{"happy path record exists", 18, http.StatusOK, `{}`},
		{"non-existent record", 445, http.StatusNotFound, `{}`},
		{"id zero test", 0, http.StatusNotFound, `{}`},
		{"id minus 1 test", -1, http.StatusNotFound, `{}`},
	}

	for _, test := range tests {
		urlUnderTest := fmt.Sprintf("%s/%d", baseUrl, test.id)
		t.Logf("Url under test : %s, test name %s", urlUnderTest, test.testName)
		req := httptest.NewRequest(http.MethodDelete, urlUnderTest, nil)
		// add on the path variable... we hope
		req.SetPathValue("id", strconv.Itoa(test.id))
		rec := httptest.NewRecorder()
		//deleteHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != test.expectedCode {
			t.Errorf("For URL %s, expected status %d, got %d", urlUnderTest, test.expectedCode, res.StatusCode)
		}

		body := rec.Body.String()
		t.Logf("returned body : %s", body)
		// body := rec.Body.String()
		// if body != test.expectedBody {
		// 	t.Errorf("For URL %s, expected body %s, got %s", test.url, test.expectedBody, body)
		// }
	}

}
