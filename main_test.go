package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestDecodeUserInput(t *testing.T) {
	t.Parallel()

	request := httptest.NewRequest(http.MethodPost, "/user/new", strings.NewReader(`{"name":" Ada ","email":"ada@example.com"}`))
	recorder := httptest.NewRecorder()

	input, ok := decodeUserInput(recorder, request)
	if !ok {
		t.Fatalf("expected valid input, got status %d", recorder.Code)
	}
	if input.Name == nil || *input.Name != " Ada " {
		t.Fatalf("unexpected name: %#v", input.Name)
	}
}

func TestDecodeUserInputRejectsUnknownAndMultipleJSONValues(t *testing.T) {
	t.Parallel()

	tests := []string{
		`{"name":"Ada","email":"ada@example.com","role":"admin"}`,
		`{"name":"Ada","email":"ada@example.com"} {}`,
	}
	for _, body := range tests {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/user/new", strings.NewReader(body))
		if _, ok := decodeUserInput(recorder, request); ok {
			t.Fatalf("expected body %q to be rejected", body)
		}
		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", recorder.Code)
		}
	}
}

func TestUserID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		value string
		valid bool
	}{
		{value: "42", valid: true},
		{value: "0", valid: false},
		{value: "abc", valid: false},
	}
	for _, test := range tests {
		recorder := httptest.NewRecorder()
		request := mux.SetURLVars(httptest.NewRequest(http.MethodGet, "/user/"+test.value, nil), map[string]string{"id": test.value})
		id, ok := userID(recorder, request)
		if ok != test.valid {
			t.Fatalf("id %q: expected valid=%t, got valid=%t", test.value, test.valid, ok)
		}
		if ok && id != 42 {
			t.Fatalf("expected ID 42, got %d", id)
		}
	}
}
