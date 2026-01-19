//
// Copyright (c) 2025 S.EE Development Team
//
// This source code is licensed under the MIT License,
// which is located in the LICENSE file in the source tree's root directory.
//
// File: client_test.go
// Author: S.EE Development Team <dev@s.ee>
// File Created: 2025-11-28 11:26:21
//
// Modified By: S.EE Development Team <dev@s.ee>
// Last Modified: 2025-12-04 17:59:00
//

package seesdk

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	baseURL := DefaultBaseURL

	if os.Getenv("SEE_API_KEY") == "" {
		t.Skip("SEE_API_KEY not set, skipping integration test")
	}

	if os.Getenv("SEE_BASE_URL") != "" {
		baseURL = os.Getenv("SEE_BASE_URL")
	}

	client := NewClient(Config{
		BaseURL: baseURL,
		APIKey:  os.Getenv("SEE_API_KEY"),
	})

	if client == nil {
		t.Fatal("Expected client to be created")
	}

	domains, err := client.GetDomains()
	if err != nil {
		t.Fatal("Expected no error, got:", err)
	}

	if len(domains.Data.Domains) == 0 {
		t.Fatal("Expected at least one domain, got zero")
	}

	tags, err := client.GetTags()
	if err != nil {
		t.Fatal("Expected no error, got:", err)
	}

	if tags.Code != 200 {
		t.Fatal("Expected response code 200, got:", tags.Code)
	}

	response, err := client.CreateShortURL(CreateShortURLRequest{
		Domain:    "a.see-test.com",
		TargetURL: "https://www.google.com/",
	})

	if err != nil {
		t.Fatal("Expected no error, got:", err)
	}

	if response.Code != 200 {
		t.Errorf("Expected response code 200, got: %d", response.Code)
	}

	result, err := client.UpdateShortURL(UpdateShortURLRequest{
		Domain:    "a.see-test.com",
		Slug:      response.Data.Slug,
		Title:     "Google",
		TargetURL: "https://www.google.com/search?q=see+sdk",
	})

	if err != nil {
		t.Fatal("Expected no error on delete, got:", err)
	}

	if result.Code != 200 {
		t.Errorf("Expected delete response code 200, got: %d", result.Code)
	}

	result2, err := client.DeleteShortURL(DeleteURLRequest{
		Domain: "a.see-test.com",
		Slug:   response.Data.Slug,
	})

	if err != nil {
		t.Fatal("Expected no error on delete, got:", err)
	}

	if result2.Code != 200 {
		t.Errorf("Expected delete response code 200, got: %d", result.Code)
	}
}

func TestTextOperations(t *testing.T) {
	baseURL := DefaultBaseURL

	if os.Getenv("SEE_API_KEY") == "" {
		t.Skip("SEE_API_KEY not set, skipping integration test")
	}

	if os.Getenv("SEE_BASE_URL") != "" {
		baseURL = os.Getenv("SEE_BASE_URL")
	}

	client := NewClient(Config{
		BaseURL: baseURL,
		APIKey:  os.Getenv("SEE_API_KEY"),
	})

	if client == nil {
		t.Fatal("Expected client to be created")
	}

	// 1. Create Text
	createResp, err := client.CreateText(CreateTextRequest{
		Domain:  "ba.sh",
		Content: "Hello, World! This is a test text.",
		Title:   "Test Text",
	})

	if err != nil {
		t.Fatal("Expected no error on create text, got:", err)
	}

	if createResp.Code != 200 {
		t.Errorf("Expected create response code 200, got: %d", createResp.Code)
	}

	if createResp.Data.Slug == "" {
		t.Fatal("Expected slug to be returned")
	}

	// 2. Update Text
	updateResp, err := client.UpdateText(UpdateTextRequest{
		Domain:  "ba.sh",
		Slug:    createResp.Data.Slug,
		Content: "Hello, World! This is an updated test text.",
		Title:   "Updated Test Text",
	})

	if err != nil {
		t.Fatal("Expected no error on update text, got:", err)
	}

	if updateResp.Code != 200 {
		t.Errorf("Expected update response code 200, got: %d", updateResp.Code)
	}

	// 3. Delete Text
	deleteResp, err := client.DeleteText(DeleteTextRequest{
		Domain: "ba.sh",
		Slug:   createResp.Data.Slug,
	})

	if err != nil {
		t.Fatal("Expected no error on delete text, got:", err)
	}

	if deleteResp.Code != 200 {
		t.Errorf("Expected delete response code 200, got: %d", deleteResp.Code)
	}
}

func TestUploadFile(t *testing.T) {
	if os.Getenv("SEE_API_KEY") == "" {
		t.Skip("SEE_API_KEY not set, skipping integration test")
	}

	baseURL := DefaultBaseURL
	if os.Getenv("SEE_BASE_URL") != "" {
		baseURL = os.Getenv("SEE_BASE_URL")
	}

	client := NewClient(Config{
		BaseURL: baseURL,
		APIKey:  os.Getenv("SEE_API_KEY"),
	})

	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "example.*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Use unique content to avoid duplicates on server side
	content := []byte(fmt.Sprintf("Hello, S.EE! %d", time.Now().UnixNano()))
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Open the file for reading
	file, err := os.Open(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	// Upload the file
	resp, err := client.UploadFile("test.txt", file)
	if err != nil {
		t.Fatal("Expected no error on upload file, got:", err)
	}

	if resp.Code != 200 {
		t.Errorf("Expected response code 200, got: %d", resp.Code)
	}

	if resp.Data.URL == "" {
		t.Error("Expected URL in response")
	}

	if resp.Data.Delete == "" {
		t.Error("Expected Delete key in response")
	}

	// Delete the file
	deleteResp, err := client.DeleteFile(resp.Data.Hash)
	if err != nil {
		t.Fatal("Expected no error on delete file, got:", err)
	}

	if !deleteResp.Success {
		t.Errorf("Expected success true, got false")
	}
}

func TestGetFileDomains(t *testing.T) {
	if os.Getenv("SEE_API_KEY") == "" {
		t.Skip("SEE_API_KEY not set, skipping integration test")
	}

	baseURL := DefaultBaseURL
	if os.Getenv("SEE_BASE_URL") != "" {
		baseURL = os.Getenv("SEE_BASE_URL")
	}

	client := NewClient(Config{
		BaseURL: baseURL,
		APIKey:  os.Getenv("SEE_API_KEY"),
	})

	domains, err := client.GetFileDomains()
	if err != nil {
		t.Fatal("Expected no error, got:", err)
	}

	if domains.Code != 0 {
		t.Errorf("Expected response code 0, got: %d", domains.Code)
	}

	if len(domains.Data.Domains) == 0 {
		t.Fatal("Expected at least one file domain, got zero")
	}
}
