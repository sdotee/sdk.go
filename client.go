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
	"time"
)

// Version is the current version of the SDK.
const Version = "1.5.0"

const DefaultBaseURL = "https://s.ee/api/v1"
const DefaultTimeout = 30 * time.Second

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

// do sets common headers, executes the request, and returns the response body.
func (c *Client) do(req *http.Request) ([]byte, error) {
	req.Header.Set("User-Agent", "see-go-sdk/"+Version)
	if c.APIKey != "" {
		req.Header.Set("Authorization", c.APIKey)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
