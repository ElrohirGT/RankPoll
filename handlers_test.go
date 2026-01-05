package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func createOrLoginUser(url string, req CreateOrLoginUserRequest) (*http.Response, error) {
	var reqBody bytes.Buffer
	err := json.NewEncoder(&reqBody).Encode(req)
	if err != nil {
		return nil, err
	}
	return http.Post(url, "application/json", &reqBody)
}

func TestCreateOrLoginUser(t *testing.T) {
	tests := []struct {
		name  string // description of this test case
		doReq func(url string, t *testing.T)
	}{
		{
			name: "Register new user",
			doReq: func(url string, t *testing.T) {
				resp, err := createOrLoginUser(url, CreateOrLoginUserRequest{Username: "fagd", Password: "12345"})
				if err != nil {
					t.Fatalf("Failed to make request: %s", err)
				}

				respBody, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("Failed to read body: %s", err)
				}

				bodyStr := strings.TrimSpace(string(respBody))
				if !strings.Contains(bodyStr, "Registered") {
					t.Fatalf("Response doesn't contain registered!")
				}
			},
		},
		{
			name: "Login user",
			doReq: func(url string, t *testing.T) {
				_, err := createOrLoginUser(url, CreateOrLoginUserRequest{Username: "fagd", Password: "12345"})
				if err != nil {
					t.Fatalf("Failed to make request: %s", err)
				}

				resp, err := createOrLoginUser(url, CreateOrLoginUserRequest{Username: "fagd", Password: "12345"})
				if err != nil {
					t.Fatalf("Failed to make request: %s", err)
				}

				respBody, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("Failed to read body: %s", err)
				}

				bodyStr := strings.TrimSpace(string(respBody))
				if !strings.Contains(bodyStr, "logging") {
					t.Fatalf("Response doesn't contain logging!")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(CreateOrLoginUser))
			defer server.Close()
			defer CleanGlobalState()

			tt.doReq(server.URL, t)
		})
	}
}

func createPoll(url string, req CreatePollRequest) (*http.Response, error) {
	var reqBody bytes.Buffer
	err := json.NewEncoder(&reqBody).Encode(req)
	if err != nil {
		return nil, err
	}
	return http.Post(url, "application/json", &reqBody)
}

func TestCreatePoll(t *testing.T) {
	tests := []struct {
		name  string // description of this test case
		doReq func(url string, t *testing.T)
	}{
		{
			name: "Create new basic poll",
			doReq: func(url string, t *testing.T) {
				resp, err := createPoll(url, CreatePollRequest{
					Title:       "Favorite Profesion?",
					PollOptions: []string{"Teacher", "Doctor", "Plumber"},
					PollUntil:   time.Now(),
				})
				if err != nil {
					t.Fatalf("Failed to make make request: %s", err)
				}

				respBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("Failed to read body of request: %s", err)
				}

				if !strings.Contains(string(respBytes), "Success") {
					t.Fatalf("Failed to create basic poll! Body: %s", respBytes)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(CreatePoll))
			defer server.Close()
			defer CleanGlobalState()

			tt.doReq(server.URL, t)
		})
	}
}
