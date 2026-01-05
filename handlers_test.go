package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateOrLoginUser(t *testing.T) {
	tests := []struct {
		name  string // description of this test case
		doReq func(url string, t *testing.T)
	}{
		{
			name: "Register new user",
			doReq: func(url string, t *testing.T) {
				resp, err := http.Post(url, "application/json", strings.NewReader(`{"Username": "fagd", "Password": "12345"}`))
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
				_, err := http.Post(url, "application/json", strings.NewReader(`{"Username": "fagd", "Password": "12345"}`))
				if err != nil {
					t.Fatalf("Failed to make request: %s", err)
				}

				resp, err := http.Post(url, "application/json", strings.NewReader(`{"Username": "fagd", "Password": "12345"}`))
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
