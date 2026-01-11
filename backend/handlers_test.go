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

func emulateHttp[T any](method string, body T, handler http.Handler) (*http.Response, error) {
	var reqBody bytes.Buffer
	err := json.NewEncoder(&reqBody).Encode(body)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(method, "/", &reqBody)
	if err != nil {
		return nil, err
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httpReq)
	return w.Result(), nil
}

func createOrLoginUser(req CreateOrLoginUserRequest) (*http.Response, error) {
	return emulateHttp(string(http.MethodPost), req, http.HandlerFunc(CreateOrLoginUser))
}

func TestCreateOrLoginUser(t *testing.T) {
	tests := []struct {
		name  string // description of this test case
		doReq func(t *testing.T)
	}{
		{
			name: "Register new user",
			doReq: func(t *testing.T) {
				resp, err := createOrLoginUser(CreateOrLoginUserRequest{Username: "fagd", Password: "12345"})
				if err != nil {
					t.Fatalf("Failed to make request: %s\n", err)
				}

				respBody, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("Failed to read body: %s\n", err)
				}

				bodyStr := strings.TrimSpace(string(respBody))
				if !strings.Contains(bodyStr, "Registered") {
					t.Fatalf("Response doesn't contain registered!\n")
				}
			},
		},
		{
			name: "Login user",
			doReq: func(t *testing.T) {
				_, err := createOrLoginUser(CreateOrLoginUserRequest{Username: "fagd", Password: "12345"})
				if err != nil {
					t.Fatalf("Failed to make request: %s\n", err)
				}

				resp, err := createOrLoginUser(CreateOrLoginUserRequest{Username: "fagd", Password: "12345"})
				if err != nil {
					t.Fatalf("Failed to make request: %s\n", err)
				}

				respBody, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("Failed to read body: %s\n", err)
				}

				bodyStr := strings.TrimSpace(string(respBody))
				if !strings.Contains(bodyStr, "logging") {
					t.Fatalf("Response doesn't contain logging!\n")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CleanGlobalState()
			tt.doReq(t)
		})
	}
}

func createPoll(req CreatePollRequest) (*http.Response, error) {
	return emulateHttp(string(http.MethodPost), req, http.HandlerFunc(CreatePoll))
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
					Title:           "Favorite Profesion?",
					PollOptions:     []string{"Teacher", "Doctor", "Plumber"},
					PollingDuration: 0,
				})
				if err != nil {
					t.Fatalf("Failed to make make request: %s\n", err)
				}

				respBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("Failed to read body of request: %s\n", err)
				}

				if !strings.Contains(string(respBytes), "Success") {
					t.Fatalf("Failed to create basic poll! Body: %s\n", respBytes)
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
	httpReq, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		return nil, err
	}
	httpReq.SetPathValue("pollId", pollId.String())

	w := httptest.NewRecorder()
	GetPollInfo(w, httpReq)
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
					Title:           "Favorite Profesion?",
					PollOptions:     []string{"Teacher", "Doctor", "Plumber"},
					PollingDuration: 0,
				})
				if err != nil {
					t.Fatalf("Failed to make make request: %s\n", err)
				}

				var createPollResponse CreatePollResponse
				err = json.NewDecoder(resp.Body).Decode(&createPollResponse)
				if err != nil {
					t.Fatalf("Failed to parse body: %s\n", err)
				}

				resp, err = getPollInfo(createPollResponse.PollId)
				if err != nil {
					t.Fatalf("Failed to make request: %s\n", err)
				}

				if resp.StatusCode != http.StatusOK {
					bodyStr, _ := io.ReadAll(resp.Body)
					t.Fatalf("The request failed somehow (%d)! %s\n", resp.StatusCode, bodyStr)
				}

				var roomInfo Room
				err = json.NewDecoder(resp.Body).Decode(&roomInfo)
				if err != nil {
					t.Fatalf("Failed to decode body: %s\n", err)
				}

				if roomInfo.Id != createPollResponse.PollId {
					t.Fatalf("Ids don't match! %s != %s\n", roomInfo.Id, createPollResponse.PollId)
				}

				if len(roomInfo.Options) != 3 {
					t.Fatalf("Room info options isn't 3!!! (%d)\n", len(roomInfo.Options))
				}
			},
		},
		{
			name: "Create and get poll info after voting",
			doReq: func(t *testing.T) {
				resp, err := createPoll(CreatePollRequest{
					Title:           "Favorite Profesion?",
					PollOptions:     []string{"Teacher", "Doctor", "Plumber"},
					PollingDuration: 300 * time.Millisecond,
				})
				if err != nil {
					t.Fatalf("Failed to make make request: %s\n", err)
				}

				var createPollResponse CreatePollResponse
				err = json.NewDecoder(resp.Body).Decode(&createPollResponse)
				if err != nil {
					t.Fatalf("Failed to parse body: %s\n", err)
				}

				_, err = createOrLoginUser(CreateOrLoginUserRequest{Username: "Tyron", Password: "12345"})
				if err != nil {
					t.Fatalf("Failed to create user: %s\n", err)
				}

				_, err = createOrLoginUser(CreateOrLoginUserRequest{Username: "Yuniqua", Password: "12345"})
				if err != nil {
					t.Fatalf("Failed to create user: %s\n", err)
				}

				_, err = createOrLoginUser(CreateOrLoginUserRequest{Username: "Pablo", Password: "12345"})
				if err != nil {
					t.Fatalf("Failed to create user: %s\n", err)
				}

				_, err = createOrLoginUser(CreateOrLoginUserRequest{Username: "Tasha", Password: "12345"})
				if err != nil {
					t.Fatalf("Failed to create user: %s\n", err)
				}

				resp, err = voteInPoll(VoteInPollRequest{Username: "Tyron", PollId: createPollResponse.PollId, Options: map[string]uint{
					"Teacher": 2,
					"Doctor":  1,
					"Plumber": 3,
				}})
				if err != nil {
					t.Fatalf("Failed to vote in poll: %s\n", err)
				}
				if resp.StatusCode != http.StatusOK {
					bodyStr, _ := io.ReadAll(resp.Body)
					t.Fatalf("Failed to vote in poll by some reason (%d): %s\n", resp.StatusCode, bodyStr)
				}

				resp, err = voteInPoll(VoteInPollRequest{Username: "Pablo", PollId: createPollResponse.PollId, Options: map[string]uint{
					"Teacher": 3,
					"Doctor":  1,
					"Plumber": 2,
				}})
				if err != nil {
					t.Fatalf("Failed to vote in poll: %s\n", err)
				}
				if resp.StatusCode != http.StatusOK {
					bodyStr, _ := io.ReadAll(resp.Body)
					t.Fatalf("Failed to vote in poll by some reason (%d): %s\n", resp.StatusCode, bodyStr)
				}

				resp, err = voteInPoll(VoteInPollRequest{Username: "Yuniqua", PollId: createPollResponse.PollId, Options: map[string]uint{
					"Teacher": 1,
					"Doctor":  2,
					"Plumber": 1,
				}})
				if err != nil {
					t.Fatalf("Failed to vote in poll: %s\n", err)
				}
				if resp.StatusCode != http.StatusOK {
					bodyStr, _ := io.ReadAll(resp.Body)
					t.Fatalf("Failed to vote in poll by some reason (%d): %s\n", resp.StatusCode, bodyStr)
				}

				resp, err = voteInPoll(VoteInPollRequest{Username: "Tasha", PollId: createPollResponse.PollId, Options: map[string]uint{
					"Teacher": 1,
					"Doctor":  2,
					"Plumber": 2,
				}})
				if err != nil {
					t.Fatalf("Failed to vote in poll: %s\n", err)
				}
				if resp.StatusCode != http.StatusOK {
					bodyStr, _ := io.ReadAll(resp.Body)
					t.Fatalf("Failed to vote in poll by some reason (%d): %s\n", resp.StatusCode, bodyStr)
				}

				time.Sleep(300 * time.Millisecond)
				resp, err = getPollInfo(createPollResponse.PollId)
				if err != nil {
					t.Fatalf("Failed to get poll info: %s\n", err)
				}

				var respBody Room
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				if err != nil {
					t.Fatalf("Can't parse body of poll info: %s\n", err)
				}

				if respBody.Summary == nil {
					t.Fatalf("Polling hasn't ended but it should have!\n")
				}

				if respBody.Summary.Winner != "Doctor" {
					t.Fatalf("Doctor didn't win! Instead %s won.\n%#v", respBody.Summary.Winner, *respBody.Summary)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CleanGlobalState()
			tt.doReq(t)
		})
	}
}

func voteInPoll(req VoteInPollRequest) (*http.Response, error) {
	return emulateHttp(http.MethodPost, req, http.HandlerFunc(VoteInPoll))
}

func TestVoteInPoll(t *testing.T) {
	tests := []struct {
		name  string // description of this test case
		doReq func(*testing.T)
	}{
		{
			name: "Vote in existent valid poll",
			doReq: func(t *testing.T) {
				_, err := createOrLoginUser(CreateOrLoginUserRequest{Username: "FAGD", Password: "12345"})
				if err != nil {
					t.Fatalf("Failed to create user: %s\n", err)
				}

				resp, err := createPoll(CreatePollRequest{
					Title:           "Lenguaje",
					PollingDuration: 5 * time.Second,
					PollOptions: []string{
						"Español", "Alemán", "Inglés",
					},
				})
				if err != nil {
					t.Fatalf("Failed to create poll: %s\n", err)
				}

				if resp.StatusCode != http.StatusOK {
					bodyStr, _ := io.ReadAll(resp.Body)
					t.Fatalf("Poll creation failed somehow (%d): %s\n", resp.StatusCode, &bodyStr)
				}

				var pollResponse CreatePollResponse
				err = json.NewDecoder(resp.Body).Decode(&pollResponse)
				if err != nil {
					t.Fatalf("Failed to decode response: %s\n", err)
				}

				resp, err = voteInPoll(VoteInPollRequest{
					Username: "FAGD",
					PollId:   pollResponse.PollId,
					Options: map[string]uint{
						"Español": 1,
						"Alemán":  3,
						"Inglés":  2,
					},
				})
				if err != nil {
					t.Fatalf("Failed make request: %s\n", err)
				}

				if resp.StatusCode != http.StatusOK {
					bodyStr, _ := io.ReadAll(resp.Body)
					t.Fatalf("Poll voting failed somehow (%d): %s\n", resp.StatusCode, &bodyStr)
				}
			},
		},
		{
			name: "Vote in existent non valid poll",
			doReq: func(t *testing.T) {
				_, err := createOrLoginUser(CreateOrLoginUserRequest{Username: "FAGD", Password: "12345"})
				if err != nil {
					t.Fatalf("Failed to create user: %s\n", err)
				}

				resp, err := createPoll(CreatePollRequest{
					Title:           "Lenguaje",
					PollingDuration: 0,
					PollOptions: []string{
						"Español", "Alemán", "Inglés",
					},
				})
				if err != nil {
					t.Fatalf("Failed to create poll: %s\n", err)
				}

				if resp.StatusCode != http.StatusOK {
					bodyStr, _ := io.ReadAll(resp.Body)
					t.Fatalf("Poll creation failed somehow (%d): %s\n", resp.StatusCode, &bodyStr)
				}

				var pollResponse CreatePollResponse
				err = json.NewDecoder(resp.Body).Decode(&pollResponse)
				if err != nil {
					t.Fatalf("Failed to decode response: %s\n", err)
				}

				resp, err = voteInPoll(VoteInPollRequest{
					Username: "FAGD",
					PollId:   pollResponse.PollId,
					Options: map[string]uint{
						"Español": 1,
						"Alemán":  3,
						"Inglés":  2,
					},
				})
				if err != nil {
					t.Fatalf("Failed to make request: %s\n", err)
				}

				if resp.StatusCode == http.StatusOK {
					bodyStr, _ := io.ReadAll(resp.Body)
					t.Fatalf("Poll voting should have failed but it didn't (%d): %s\n", resp.StatusCode, &bodyStr)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CleanGlobalState()
			tt.doReq(t)
		})
	}
}
