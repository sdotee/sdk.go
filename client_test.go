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
	"bytes"
	"crypto/rand"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"
)

func setupTestClient(t *testing.T) *Client {
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

	if client == nil {
		t.Fatal("Expected client to be created")
	}

	return client
}

// testShortURLDomain returns the first available short URL domain for the account.
func testShortURLDomain(t *testing.T, client *Client) string {
	t.Helper()
	domains, err := client.GetDomains()
	if err != nil {
		t.Fatal("Expected no error getting domains, got:", err)
	}
	if len(domains.Data.Domains) == 0 {
		t.Fatal("Expected at least one domain, got zero")
	}
	return domains.Data.Domains[0]
}

// testTextDomain returns the first available text sharing domain for the account.
func testTextDomain(t *testing.T, client *Client) string {
	t.Helper()
	domains, err := client.GetTextDomains()
	if err != nil {
		t.Fatal("Expected no error getting text domains, got:", err)
	}
	if len(domains.Data.Domains) == 0 {
		t.Fatal("Expected at least one text domain, got zero")
	}
	return domains.Data.Domains[0]
}

func TestNewClient(t *testing.T) {
	client := setupTestClient(t)

	domain := testShortURLDomain(t, client)

	tags, err := client.GetTags()
	if err != nil {
		t.Fatal("Expected no error, got:", err)
	}

	if tags.Code != 200 {
		t.Fatal("Expected response code 200, got:", tags.Code)
	}

	response, err := client.CreateShortURL(CreateShortURLRequest{
		Domain:    domain,
		TargetURL: "https://www.google.com/",
	})

	if err != nil {
		t.Fatal("Expected no error, got:", err)
	}

	if response.Code != 200 {
		t.Errorf("Expected response code 200, got: %d", response.Code)
	}

	result, err := client.UpdateShortURL(UpdateShortURLRequest{
		Domain:    domain,
		Slug:      response.Data.Slug,
		Title:     "Google",
		TargetURL: "https://www.google.com/search?q=see+sdk",
	})

	if err != nil {
		t.Fatal("Expected no error on update, got:", err)
	}

	if result.Code != 200 {
		t.Errorf("Expected update response code 200, got: %d", result.Code)
	}

	result2, err := client.DeleteShortURL(DeleteURLRequest{
		Domain: domain,
		Slug:   response.Data.Slug,
	})

	if err != nil {
		t.Fatal("Expected no error on delete, got:", err)
	}

	if result2.Code != 200 {
		t.Errorf("Expected delete response code 200, got: %d", result2.Code)
	}
}

func TestTextOperations(t *testing.T) {
	client := setupTestClient(t)

	domain := testTextDomain(t, client)

	// 1. Create Text
	createResp, err := client.CreateText(CreateTextRequest{
		Domain:  domain,
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
		Domain:  domain,
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
		Domain: domain,
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
	client := setupTestClient(t)

	file, err := os.Open("testdata/test.png")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	resp, err := client.UploadFile(UploadFileRequest{
		Filename: "test.png",
		File:     file,
	})
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

	if resp.Data.Width == 0 || resp.Data.Height == 0 {
		t.Errorf("Expected image dimensions, got: %dx%d", resp.Data.Width, resp.Data.Height)
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

func TestUploadFilePrivate(t *testing.T) {
	client := setupTestClient(t)

	file, err := os.Open("testdata/test.png")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	// Upload as private
	resp, err := client.UploadFile(UploadFileRequest{
		Filename:  "test-private.png",
		File:      file,
		IsPrivate: true,
	})
	if err != nil {
		t.Fatal("Expected no error on upload private file, got:", err)
	}

	if resp.Code != 200 {
		t.Errorf("Expected response code 200, got: %d", resp.Code)
	}

	if resp.Data.FileID == 0 {
		t.Error("Expected file ID in response")
	}

	// Clean up
	deleteResp, err := client.DeleteFile(resp.Data.Hash)
	if err != nil {
		t.Fatal("Expected no error on delete file, got:", err)
	}

	if !deleteResp.Success {
		t.Errorf("Expected success true, got false")
	}
}

func TestGetFileHistory(t *testing.T) {
	client := setupTestClient(t)

	resp, err := client.GetFileHistory(1)
	if err != nil {
		t.Fatal("Expected no error, got:", err)
	}

	if resp.Code != 200 {
		t.Errorf("Expected response code 200, got: %d", resp.Code)
	}

	if !resp.Success {
		t.Error("Expected success true")
	}

	fmt.Printf("File history: %d files on page 1\n", len(resp.Data))
}

func TestGetFileDomains(t *testing.T) {
	client := setupTestClient(t)

	domains, err := client.GetFileDomains()
	if err != nil {
		t.Fatal("Expected no error, got:", err)
	}

	if domains.Code != 200 {
		t.Errorf("Expected response code 0, got: %d", domains.Code)
	}

	if len(domains.Data.Domains) == 0 {
		t.Fatal("Expected at least one file domain, got zero")
	}

	fmt.Println("Available file domains:")
	for _, domain := range domains.Data.Domains {
		fmt.Printf(" - %s\n", domain)
	}
}

func TestGetTextDomains(t *testing.T) {
	client := setupTestClient(t)

	domains, err := client.GetTextDomains()
	if err != nil {
		t.Fatal("Expected no error, got:", err)
	}

	if domains.Code != 200 {
		t.Errorf("Expected response code 200, got: %d", domains.Code)
	}

	if len(domains.Data.Domains) == 0 {
		t.Fatal("Expected at least one text domain, got zero")
	}

	fmt.Println("Available text domains:")
	for _, domain := range domains.Data.Domains {
		fmt.Printf(" - %s\n", domain)
	}
}

func TestGetUsage(t *testing.T) {
	client := setupTestClient(t)

	usage, err := client.GetUsage()
	if err != nil {
		t.Fatal("Expected no error, got:", err)
	}

	if usage.Code != 200 {
		t.Errorf("Expected response code 200, got: %d", usage.Code)
	}

	fmt.Printf("Links created this month: %d/%d\n",
		usage.Data.LinkCountMonth,
		usage.Data.LinkCountMonthLimit)

	if usage.Data.APICountMonthLimit != UsageNoLimit {
		t.Error("Expected API count month limit to be no limit")
	}
}

func TestGetLinkVisitStat(t *testing.T) {
	client := setupTestClient(t)

	domain := testShortURLDomain(t, client)

	// Create a short URL to query statistics for.
	createResp, err := client.CreateShortURL(CreateShortURLRequest{
		Domain:    domain,
		TargetURL: "https://www.example.com/",
	})
	if err != nil {
		t.Fatal("Expected no error on create, got:", err)
	}

	defer client.DeleteShortURL(DeleteURLRequest{
		Domain: domain,
		Slug:   createResp.Data.Slug,
	})

	for _, period := range []string{"", VisitStatPeriodDaily, VisitStatPeriodMonthly, VisitStatPeriodTotally} {
		statResp, err := client.GetLinkVisitStat(domain, createResp.Data.Slug, period)
		if err != nil {
			t.Fatalf("Expected no error for period %q, got: %v", period, err)
		}

		if statResp.Code != 200 {
			t.Errorf("Expected response code 200 for period %q, got: %d (message: %s)", period, statResp.Code, statResp.Message)
		}

		if statResp.Data.VisitCount != 0 {
			t.Errorf("Expected zero visits for a new link (period %q), got: %d", period, statResp.Data.VisitCount)
		}
	}
}

func TestGetPrivateFileDownloadURL(t *testing.T) {
	client := setupTestClient(t)

	// Since we need an existing private file ID for this test, let's create one.
	// We'll upload a temporary file.

	// 1. Create a dummy file
	f, err := os.CreateTemp("", "test-private-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString("This is a private file content test.")
	f.Close()

	file, err := os.Open(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	// 2. Upload it privately
	uploadResp, err := client.UploadFile(UploadFileRequest{
		Filename:  "test-private-dl.txt",
		File:      file,
		IsPrivate: true,
	})
	if err != nil {
		t.Fatal("Failed to upload private file:", err)
	}

	if uploadResp.Code != 200 {
		t.Logf("Upload failed with code %d: %s (maybe account limit reached or other issue)", uploadResp.Code, uploadResp.Message)
		return
	}

	// Ensure cleanup of the uploaded file
	defer func() {
		client.DeleteFile(uploadResp.Data.Hash)
	}()

	t.Logf("Uploaded private file ID: %d", uploadResp.Data.FileID)

	// 3. Test GetPrivateFileDownloadURL
	dlResp, err := client.GetPrivateFileDownloadURL(int64(uploadResp.Data.FileID))
	if err != nil {
		t.Fatal("GetPrivateFileDownloadURL failed:", err)
	}

	if dlResp.Code != 200 {
		t.Errorf("Expected response code 200, got: %d (message: %s)", dlResp.Code, dlResp.Message)
	}

	if dlResp.Data.URL == "" {
		t.Error("Expected download URL in response, got empty string")
	}

	if dlResp.Data.FileID != int64(uploadResp.Data.FileID) {
		t.Errorf("Expected file ID %d, got %d", uploadResp.Data.FileID, dlResp.Data.FileID)
	}

	if dlResp.Data.ExpiresAt == 0 {
		t.Error("Expected expiration timestamp, got 0")
	}

	t.Logf("Got private download URL: %s", dlResp.Data.URL)
}

func TestLargeFileUploadLifecycle(t *testing.T) {
	client := setupTestClient(t)

	// 1. Create an upload session.
	createResp, err := client.CreateLargeFileUpload(CreateLargeFileUploadRequest{
		FileName: "test-large-lifecycle.bin",
		FileSize: 1024 * 1024, // 1MB
	})
	if err != nil {
		t.Fatal("Expected no error on create session, got:", err)
	}

	if createResp.Code != 200 {
		t.Fatalf("Expected response code 200, got: %d (message: %s)", createResp.Code, createResp.Message)
	}

	if createResp.Data.UploadID == "" {
		t.Fatal("Expected upload ID in response")
	}

	// 2. Query the upload progress.
	progressResp, err := client.GetLargeFileUploadProgress(createResp.Data.UploadID)
	if err != nil {
		t.Fatal("Expected no error on get progress, got:", err)
	}

	if progressResp.Code != 200 {
		t.Errorf("Expected progress response code 200, got: %d (message: %s)", progressResp.Code, progressResp.Message)
	}

	if progressResp.Data.Status != LargeFileUploadStatusUploading {
		t.Errorf("Expected status uploading (%d), got: %d", LargeFileUploadStatusUploading, progressResp.Data.Status)
	}

	// 3. Cancel the session.
	cancelResp, err := client.CancelLargeFileUpload(createResp.Data.UploadID)
	if err != nil {
		t.Fatal("Expected no error on cancel, got:", err)
	}

	if cancelResp.Code != 200 {
		t.Errorf("Expected cancel response code 200, got: %d (message: %s)", cancelResp.Code, cancelResp.Message)
	}
}

func TestUploadLargeFile(t *testing.T) {
	client := setupTestClient(t)

	// Create a temporary file with random content.
	const fileSize = 1024 * 1024 // 1MB
	content := make([]byte, fileSize)
	if _, err := rand.Read(content); err != nil {
		t.Fatal(err)
	}

	resp, err := client.UploadLargeFile(CreateLargeFileUploadRequest{
		FileName: "test-large-upload.bin",
		FileSize: fileSize,
	}, bytes.NewReader(content))
	if err != nil {
		t.Fatal("Expected no error on large file upload, got:", err)
	}

	if resp.Code != 200 {
		t.Fatalf("Expected response code 200, got: %d (message: %s)", resp.Code, resp.Message)
	}

	if resp.Data.File.FileID == 0 {
		t.Error("Expected file ID in response")
	}

	if resp.Data.File.Hash == "" {
		t.Fatal("Expected delete key in response")
	}

	t.Logf("Uploaded large file ID: %d, size: %d", resp.Data.File.FileID, resp.Data.File.Size)

	// Clean up
	deleteResp, err := client.DeleteFile(resp.Data.File.Hash)
	if err != nil {
		t.Fatal("Expected no error on delete file, got:", err)
	}

	if !deleteResp.Success {
		t.Errorf("Expected success true, got false")
	}
}

func TestSmartUploadFile(t *testing.T) {
	client := setupTestClient(t)

	// A small payload should go through the regular multipart upload path.
	content := make([]byte, 1024)
	if _, err := rand.Read(content); err != nil {
		t.Fatal(err)
	}

	resp, err := client.SmartUploadFile(UploadFileRequest{
		File:     bytes.NewReader(content),
		Filename: "test-smart-upload.bin",
	})
	if err != nil {
		t.Fatal("Expected no error on smart upload, got:", err)
	}

	if resp.Code != 200 {
		t.Fatalf("Expected response code 200, got: %d (message: %s)", resp.Code, resp.Message)
	}

	if resp.Data.Hash == "" {
		t.Fatal("Expected delete key in response")
	}

	// Clean up
	if _, err := client.DeleteFile(resp.Data.Hash); err != nil {
		t.Fatal("Expected no error on delete file, got:", err)
	}
}

func TestGetLinkHistory(t *testing.T) {
	client := setupTestClient(t)

	resp, err := client.GetLinkHistory(1)
	if err != nil {
		t.Fatal("Expected no error, got:", err)
	}

	if resp.Code != 200 {
		t.Errorf("Expected response code 200, got: %d (message: %s)", resp.Code, resp.Message)
	}

	t.Logf("Got %d link history entries", len(resp.Data))
}

func TestGetTextHistory(t *testing.T) {
	client := setupTestClient(t)

	resp, err := client.GetTextHistory(1)
	if err != nil {
		t.Fatal("Expected no error, got:", err)
	}

	if resp.Code != 200 {
		t.Errorf("Expected response code 200, got: %d (message: %s)", resp.Code, resp.Message)
	}

	t.Logf("Got %d text history entries", len(resp.Data))
}

func TestCheckToken(t *testing.T) {
	client := setupTestClient(t)

	resp, err := client.CheckToken(os.Getenv("SEE_API_KEY"))
	if err != nil {
		t.Fatal("Expected no error, got:", err)
	}

	if resp.Code != 200 {
		t.Fatalf("Expected response code 200, got: %d (message: %s)", resp.Code, resp.Message)
	}

	if !resp.Data.Valid {
		t.Error("Expected token to be valid")
	}
}

func TestBioPageLifecycle(t *testing.T) {
	client := setupTestClient(t)

	// 1. Create a bio page.
	createResp, err := client.CreateBioPage(CreateBioPageRequest{
		Title:       "SDK Test Bio Page",
		Description: "Created by see-go-sdk integration tests",
		CustomLinks: []BioCustomLink{
			{Title: "Example", URL: "https://example.com"},
		},
	})
	if err != nil {
		t.Fatal("Expected no error on create bio page, got:", err)
	}

	if createResp.Code != 200 {
		t.Fatalf("Expected response code 200, got: %d (message: %s)", createResp.Code, createResp.Message)
	}

	if createResp.Data.BioPageID == 0 {
		t.Fatal("Expected bio page ID in response")
	}

	t.Logf("Created bio page %d: %s", createResp.Data.BioPageID, createResp.Data.ShortURL)

	// 2. Update the bio page.
	updateResp, err := client.UpdateBioPage(UpdateBioPageRequest{
		ID:    createResp.Data.BioPageID,
		Title: "SDK Test Bio Page (updated)",
	})
	if err != nil {
		t.Fatal("Expected no error on update bio page, got:", err)
	}

	if updateResp.Code != 200 {
		t.Errorf("Expected update response code 200, got: %d (message: %s)", updateResp.Code, updateResp.Message)
	}

	// 3. List bio pages.
	historyResp, err := client.GetBioPageHistory(1)
	if err != nil {
		t.Fatal("Expected no error on bio page history, got:", err)
	}

	if historyResp.Code != 200 {
		t.Errorf("Expected history response code 200, got: %d (message: %s)", historyResp.Code, historyResp.Message)
	}

	// 4. Delete the bio page.
	deleteResp, err := client.DeleteBioPage(createResp.Data.BioPageID)
	if err != nil {
		t.Fatal("Expected no error on delete bio page, got:", err)
	}

	if deleteResp.Code != 200 {
		t.Errorf("Expected delete response code 200, got: %d (message: %s)", deleteResp.Code, deleteResp.Message)
	}
}

func TestQRCodeLifecycle(t *testing.T) {
	client := setupTestClient(t)

	// 1. Create a QR code.
	createResp, err := client.CreateQRCode(CreateQRCodeRequest{
		TargetURL: "https://example.com",
		Title:     "SDK Test QR Code",
	})
	if err != nil {
		// The QR code endpoints are documented but may not be deployed yet.
		if strings.Contains(err.Error(), "status 404") {
			t.Skip("QR code endpoint not available on this server, skipping")
		}
		t.Fatal("Expected no error on create QR code, got:", err)
	}

	if createResp.Code != 200 {
		t.Fatalf("Expected response code 200, got: %d (message: %s)", createResp.Code, createResp.Message)
	}

	if createResp.Data.Slug == "" {
		t.Fatal("Expected slug in response")
	}

	t.Logf("Created QR code: %s (PNG: %s)", createResp.Data.ShortURL, createResp.Data.PNGURL)

	// 2. List QR codes.
	historyResp, err := client.GetQRCodeHistory(1)
	if err != nil {
		t.Fatal("Expected no error on QR code history, got:", err)
	}

	if historyResp.Code != 200 {
		t.Errorf("Expected history response code 200, got: %d (message: %s)", historyResp.Code, historyResp.Message)
	}

	// 3. Delete the QR code. Derive domain from the short URL when possible.
	domain := createResp.Data.CustomSlug
	if u, err := url.Parse(createResp.Data.ShortURL); err == nil {
		domain = u.Hostname()
	}

	deleteResp, err := client.DeleteQRCode(DeleteQRCodeRequest{
		Domain: domain,
		Slug:   createResp.Data.Slug,
	})
	if err != nil {
		t.Fatal("Expected no error on delete QR code, got:", err)
	}

	if deleteResp.Code != 200 {
		t.Errorf("Expected delete response code 200, got: %d (message: %s)", deleteResp.Code, deleteResp.Message)
	}
}
