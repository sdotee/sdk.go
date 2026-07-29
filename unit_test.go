//
// Copyright (c) 2025-2026 S.EE Development Team
//
// This source code is licensed under the MIT License,
// which is located in the LICENSE file in the source tree's root directory.
//
// File: unit_test.go
// Author: S.EE Development Team <dev@s.ee>
// File Created: 2026-07-29 12:00:00
//

package seesdk

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAPIErrorUsesJSONMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"code":401,"message":"invalid API key"}`))
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL})
	_, err := client.GetDomains()
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized || apiErr.Message != "invalid API key" {
		t.Fatalf("unexpected API error: %#v", apiErr)
	}
	if apiErr.Method != http.MethodGet || apiErr.Endpoint != "/domains" {
		t.Fatalf("missing request context: %#v", apiErr)
	}
}

func TestAPIErrorSummarizesHTML(t *testing.T) {
	html := "<!doctype html><html><body>upstream details</body></html>"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(html))
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL})
	_, err := client.GetDomains()
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "server returned HTML instead of the expected JSON response") {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(err.Error(), "upstream details") {
		t.Fatalf("error leaked HTML body: %v", err)
	}
}

func TestSuccessfulNonJSONResponse(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{name: "HTML", body: "  <html><body>OK</body></html>", want: "expected JSON but received HTML"},
		{name: "empty", body: "  \n", want: "expected JSON but received an empty response"},
		{name: "invalid JSON", body: "not-json", want: "invalid JSON"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(test.body))
			}))
			defer server.Close()

			client := NewClient(Config{BaseURL: server.URL})
			_, err := client.GetDomains()
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("expected error containing %q, got %v", test.want, err)
			}
		})
	}
}

func TestTUSErrorUsesSharedResponseHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("<html><body>maintenance</body></html>"))
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL})
	_, err := client.UploadLargeFileChunk("upload-id", 0, []byte("chunk"))
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusServiceUnavailable ||
		apiErr.Message != "server returned HTML instead of the expected JSON response" {
		t.Fatalf("unexpected TUS API error: %#v", apiErr)
	}
}
