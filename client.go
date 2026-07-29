//
// Copyright (c) 2025 S.EE Development Team
//
// This source code is licensed under the MIT License,
// which is located in the LICENSE file in the source tree's root directory.
//
// File: client.go
// Author: S.EE Development Team <dev@s.ee>
// File Created: 2025-11-28 11:21:45
//
// Modified By: S.EE Development Team <dev@s.ee>
// Last Modified: 2026-07-27 12:00:00
//

package seesdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

// Version is the current version of the SDK.
const Version = "1.5.0"

const DefaultBaseURL = "https://s.ee/api/v1"
const DefaultTimeout = 30 * time.Second

const maxErrorResponseSize = 64 * 1024

// userAgent is the User-Agent header value sent with every request.
const userAgent = "see-go-sdk/" + Version

// Client represents the SEE SDK client for short URL operations
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// Config contains configuration options for the Client
type Config struct {
	BaseURL string
	APIKey  string
	Timeout time.Duration
}

// APIError describes a non-successful HTTP response from the API.
type APIError struct {
	Method     string
	Endpoint   string
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API request %s %s failed (status %d): %s", e.Method, e.Endpoint, e.StatusCode, e.Message)
}

// NewClient creates a new SEE SDK client with the given configuration.
func NewClient(config Config) *Client {
	if config.BaseURL == "" {
		config.BaseURL = DefaultBaseURL
	}
	if config.Timeout == 0 {
		config.Timeout = DefaultTimeout
	}

	return &Client{
		BaseURL: config.BaseURL,
		APIKey:  config.APIKey,
		HTTPClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// doRequest executes a JSON HTTP request and returns the response body.
func (c *Client) doRequest(method, endpoint string, body any) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.BaseURL+endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	return c.do(req)
}

// doMultipartRequest executes a multipart HTTP request.
func (c *Client) doMultipartRequest(endpoint string, fieldName, filename string, r io.Reader, fields map[string]string) ([]byte, error) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()
		part, err := writer.CreateFormFile(fieldName, filename)
		if err != nil {
			_ = pw.CloseWithError(fmt.Errorf("create form file: %w", err))
			return
		}
		if _, err := io.Copy(part, r); err != nil {
			_ = pw.CloseWithError(fmt.Errorf("copy file content: %w", err))
			return
		}
		for key, value := range fields {
			if err := writer.WriteField(key, value); err != nil {
				_ = pw.CloseWithError(fmt.Errorf("write field %s: %w", key, err))
				return
			}
		}
		if err := writer.Close(); err != nil {
			_ = pw.CloseWithError(fmt.Errorf("close writer: %w", err))
		}
	}()

	req, err := http.NewRequest("POST", c.BaseURL+endpoint, pr)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return c.do(req)
}

// setCommonHeaders sets the headers shared by every API request.
func (c *Client) setCommonHeaders(req *http.Request) {
	req.Header.Set("User-Agent", userAgent)
	if c.APIKey != "" {
		req.Header.Set("Authorization", c.APIKey)
	}
}

// do sets common headers, executes the request, and returns the response body.
func (c *Client) do(req *http.Request) ([]byte, error) {
	c.setCommonHeaders(req)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	var responseReader io.Reader = resp.Body
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		responseReader = io.LimitReader(resp.Body, maxErrorResponseSize)
	}
	respBody, err := io.ReadAll(responseReader)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &APIError{
			Method:     req.Method,
			Endpoint:   req.URL.RequestURI(),
			StatusCode: resp.StatusCode,
			Message:    responseErrorMessage(resp.Header.Get("Content-Type"), respBody),
		}
	}

	return respBody, nil
}

func responseErrorMessage(contentType string, body []byte) string {
	var payload struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}
	if json.Unmarshal(body, &payload) == nil {
		if payload.Message != "" {
			return payload.Message
		}
		if payload.Error != "" {
			return payload.Error
		}
	}

	trimmed := strings.TrimSpace(string(body))
	lowerBody := strings.ToLower(trimmed)
	if strings.Contains(strings.ToLower(contentType), "text/html") ||
		strings.HasPrefix(lowerBody, "<!doctype html") || strings.HasPrefix(lowerBody, "<html") {
		return "server returned HTML instead of the expected JSON response"
	}
	if trimmed == "" {
		return "server returned an empty response"
	}

	message := strings.Join(strings.Fields(trimmed), " ")
	const maxErrorMessageLength = 200
	if len(message) > maxErrorMessageLength {
		message = message[:maxErrorMessageLength] + "..."
	}
	return message
}
