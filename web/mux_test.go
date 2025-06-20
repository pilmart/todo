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

	t.Skip("Skipping this test for now")
	baseUrl := "/todo"
	tests := []struct {
		testName     string
		id           int
		expectedCode int
		expectedBody string
	}{
		{"happy path record present", 18, http.StatusOK, `{"id":18,"description":"this should work","status":"NOT STARTED"}`},
		{"non-existent record", 445, http.StatusNotFound, `{}`},
		{"id zero test", 0, http.StatusNotFound, `{}`},
		{"id minus one test", -1, http.StatusNotFound, `{}`},
	}

	for _, test := range tests {
		urlUnderTest := fmt.Sprintf("%s/%d", baseUrl, test.id)
		t.Logf("Url under test : %s, test name %s", urlUnderTest, test.testName)
		req := httptest.NewRequest(http.MethodGet, urlUnderTest, nil)
		// add on the path variable... we hope
		req.SetPathValue("id", strconv.Itoa(test.id))
		rec := httptest.NewRecorder()
		//getHandler(rec, req)

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

// test endpoint POST /todo - need a current record
func TestPost(t *testing.T) {
	t.Skip("Skipping this test for now")
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
		urlUnderTest := baseUrl
		t.Logf("Url under test : %s, test name %s", urlUnderTest, test.testName)
		reqBody := strings.NewReader(test.expectedBody)
		req := httptest.NewRequest(http.MethodPost, urlUnderTest, reqBody)
		_ = req
		rec := httptest.NewRecorder()
		//createHandler(rec, req)

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

// test endpoint PUT /todo - need a current record
func TestPut(t *testing.T) {

	baseUrl := "/todo"
	tests := []struct {
		testName     string
		expectedCode int
		expectedBody string
	}{
		{"happy path", http.StatusOK, `{"id":21, "description":"Updated via update test","status": "NOT STARTED"}`},
		{"missing description", http.StatusOK, `{"id":21,"description":"","status": "NOT STARTED"}`},
		{"Incorrect status", http.StatusBadRequest, `{"id":21,"description":"Updated via test","status": "Incorrect"}`},
		{"empty body", http.StatusNotFound, `{}`},
		{"non-existent id", http.StatusNotFound, `{"id":99999,"description":"Updated via update test","status": "NOT STARTED"}`},
	}

	for _, test := range tests {
		urlUnderTest := baseUrl
		t.Logf("Url under test : %s, test name %s", urlUnderTest, test.testName)
		reqBody := strings.NewReader(test.expectedBody)
		req := httptest.NewRequest(http.MethodPost, urlUnderTest, reqBody)
		_ = req
		rec := httptest.NewRecorder()
		//updateHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != test.expectedCode {
			t.Errorf("Test name %s, For URL %s, expected status %d, got %d", test.testName, urlUnderTest, test.expectedCode, res.StatusCode)
		}

		body := rec.Body.String()
		t.Logf("returned body : %s", body)
		// body := rec.Body.String()
		// if body != test.expectedBody {
		// 	t.Errorf("For URL %s, expected body %s, got %s", test.url, test.expectedBody, body)
		// }
	}

}

// test endpoint DELETE /todo/{id}
func TestDelete(t *testing.T) {
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
