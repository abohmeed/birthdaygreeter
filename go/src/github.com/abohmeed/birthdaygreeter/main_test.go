package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func updateUser(days int) int {
	inputDate := time.Date(1979, time.Now().Month(), time.Now().Day()+days, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	var jsonStr = []byte(`{"dateOfBirth":"` + inputDate + `"}`)
	req, _ := http.NewRequest("PUT", "/hello/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()
	newServer().ServeHTTP(res, req)
	return res.Code
}
func TestHandleQueryBirthdate(t *testing.T) {
	//Create/update user
	updateUser(10)
	req, _ := http.NewRequest("GET", "/hello/user", nil)
	res := httptest.NewRecorder()
	newServer().ServeHTTP(res, req)
	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body is what we expect.
	expected := `{"message":"Hello, user"}` + "\n"
	if res.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: want %v got %v", expected, res.Body.String())
	}
}
func TestHandleUpdateBirthdate(t *testing.T) {
	//Create/update user
	rc := updateUser(5)
	expected := 204
	if rc != expected {
		t.Errorf("Handler returned unexpected status code: want %v got %v", expected, rc)
	}
}
func TestUpcomingBirthdate(t *testing.T) {
	req, _ := http.NewRequest("GET", "/hello/user", nil)
	res := httptest.NewRecorder()
	newServer().ServeHTTP(res, req)
	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body is what we expect.
	expected := `{"message":"Hello, user! Your birthday is in 5 days"}` + "\n"
	if res.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: want %v got %v", expected, res.Body.String())
	}
}

func TestTodayisBirthdate(t *testing.T) {
	updateUser(0)
	req, _ := http.NewRequest("GET", "/hello/user", nil)
	res := httptest.NewRecorder()
	newServer().ServeHTTP(res, req)
	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body is what we expect.
	expected := `{"message":"Hello, user! Happy birthday"}` + "\n"
	if res.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: want %v got %v", expected, res.Body.String())
	}
}
func TestInvalidJSONInput(t *testing.T) {
	var jsonStr = []byte(`{"dateOfBirth":"October 25th, 1979"}`)
	req, _ := http.NewRequest("PUT", "/hello/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()
	newServer().ServeHTTP(res, req)
	expected := 422
	if res.Code != expected {
		t.Errorf("Handler returned unexpected body: want %v got %v", expected, res.Code)
	}
	expectedMsg := `{"message":"Invalid message format, should be \"dateOfBirth\":\"yyyy-mm-dd\""}` + "\n"
	if res.Body.String() != expectedMsg {
		t.Errorf("Handler returned unexpected body: want %v got %v", expectedMsg, res.Body.String())
	}
}
func TestInvalidInput(t *testing.T) {
	var jsonStr = []byte(`This is not JSON`)
	req, _ := http.NewRequest("PUT", "/hello/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()
	newServer().ServeHTTP(res, req)
	expected := 422
	if res.Code != expected {
		t.Errorf("Handler returned unexpected body: want %v got %v", expected, res.Code)
	}
}
