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

	"github.com/google/uuid"
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

func createPoll(req CreatePollRequest) (*http.Response, error) {
	var reqBody bytes.Buffer
	err := json.NewEncoder(&reqBody).Encode(req)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequest(http.MethodPost, "/", &reqBody)
	if err != nil {
		return nil, err
	}

	w := httptest.NewRecorder()
	CreatePoll(w, httpReq)

	return w.Result(), nil
}

func TestCreatePoll(t *testing.T) {
	tests := []struct {
		name  string // description of this test case
		doReq func(t *testing.T)
	}{
		{
			name: "Create new basic poll",
			doReq: func(t *testing.T) {
				resp, err := createPoll(CreatePollRequest{
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
			tt.doReq(t)
			CleanGlobalState()
		})
	}
}

func getPollInfo(pollId uuid.UUID) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, "/", nil)
	if err != nil {
		return nil, err
	}
	req.SetPathValue("pollId", pollId.String())

	w := httptest.NewRecorder()
	GetPollInfo(w, req)

	return w.Result(), nil
}

func TestGetPollInfo(t *testing.T) {
	tests := []struct {
		name  string
		doReq func(t *testing.T)
	}{
		{
			name: "Create and get poll info",
			doReq: func(t *testing.T) {
				resp, err := createPoll(CreatePollRequest{
					Title:       "Favorite Profesion?",
					PollOptions: []string{"Teacher", "Doctor", "Plumber"},
					PollUntil:   time.Now(),
				})
				if err != nil {
					t.Fatalf("Failed to make make request: %s", err)
				}

				var createPollResponse CreatePollResponse
				err = json.NewDecoder(resp.Body).Decode(&createPollResponse)
				if err != nil {
					t.Fatalf("Failed to parse body: %s", err)
				}

				resp, err = getPollInfo(createPollResponse.PollId)
				if err != nil {
					t.Fatalf("Failed to make request: %s", err)
				}

				if resp.StatusCode != http.StatusOK {
					bodyStr, _ := io.ReadAll(resp.Body)
					t.Fatalf("The request failed somehow (%d)! %s", resp.StatusCode, bodyStr)
				}

				var roomInfo Room
				err = json.NewDecoder(resp.Body).Decode(&roomInfo)
				if err != nil {
					t.Fatalf("Failed to decode body: %s", err)
				}

				if roomInfo.Id != createPollResponse.PollId {
					t.Fatalf("Ids don't match! %s != %s", roomInfo.Id, createPollResponse.PollId)
				}

				if len(roomInfo.Options) != 3 {
					t.Fatalf("Room info options isn't 3!!! (%d)", len(roomInfo.Options))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.doReq(t)
			CleanGlobalState()
		})
	}
}
