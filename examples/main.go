//
// Copyright (c) 2025-2026 S.EE Development Team
//
// This source code is licensed under the MIT License,
// which is located in the LICENSE file in the source tree's root directory.
//
// File: main.go
// Author: S.EE Development Team <dev@s.ee>
// File Created: 2025-11-28 11:26:23
//
// Modified By: S.EE Development Team <dev@s.ee>
// Last Modified: 2025-12-04 17:58:16
//

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	seesdk "github.com/sdotee/sdk.go"
)

var (
	client        *seesdk.Client
	defaultDomain = "example.com"
)

func init() {
	baseURL := os.Getenv("SEE_API_URL")
	if baseURL == "" {
		baseURL = seesdk.DefaultBaseURL
	}

	apiKey := os.Getenv("SEE_API_KEY")
	if apiKey == "" {
		// Use a placeholder if not set, but warn the user
		log.Println("Note: SEE_API_KEY environment variable not set. Using placeholder API key.")
		apiKey = "your-api-key-here"
	}

	client = seesdk.NewClient(seesdk.Config{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Timeout: 30 * time.Second,
	})
}

func main() {
	// 1. Get Domains (and pick one for subsequent tests)
	if domain := getDomains(); domain != "" {
		defaultDomain = domain
		fmt.Printf("Using domain '%s' for tests\n", defaultDomain)
	}

	// 2. Get Tags
	getTags()

	// 3. Short URL Lifecycle: Create -> Update -> Delete
	fmt.Println("\n--- Short URL Lifecycle ---")
	if slug := createShortURL(); slug != "" {
		updateShortURL(slug)
		deleteShortURL(slug)
	}

	// 4. Custom Short URL
	fmt.Println("\n--- Custom Short URL ---")
	createCustomShortURL()

	// 5. Protected Short URL
	fmt.Println("\n--- Protected Short URL ---")
	createProtectedShortURL()

	// 6. File Operations
	fmt.Println("\n--- File Operations ---")
	fileOperations()
}

func getDomains() string {
	fmt.Println("=== Get Domains ===")
	domainsResp, err := client.GetDomains()
	if err != nil {
		log.Printf("Failed to get domains: %v\n", err)
		return ""
	}
	fmt.Printf("Available domains: %v\n\n", domainsResp.Data.Domains)

	if len(domainsResp.Data.Domains) > 0 {
		return domainsResp.Data.Domains[0]
	}
	return ""
}

func getTags() {
	fmt.Println("=== Get Tags ===")
	tagsResp, err := client.GetTags()
	if err != nil {
		log.Printf("Failed to get tags: %v\n", err)
		return
	}
	fmt.Printf("Available tags:\n")
	for _, tag := range tagsResp.Data.Tags {
		fmt.Printf("  - ID: %d, Name: %s\n", tag.ID, tag.Name)
	}
	fmt.Println()
}

func createShortURL() string {
	fmt.Println("=== Create Short URL ===")
	createResp, err := client.CreateShortURL(seesdk.CreateShortURLRequest{
		TargetURL: "https://www.example.com/very/long/url/path",
		Domain:    defaultDomain,
		Title:     "Example Link",
	})
	if err != nil {
		log.Printf("Failed to create short URL: %v\n", err)
		return ""
	}
	fmt.Printf("Response Code: %d\n", createResp.Code)
	fmt.Printf("Message: %s\n", createResp.Message)
	fmt.Printf("Slug: %s\n", createResp.Data.Slug)
	fmt.Printf("Short URL: %s\n\n", createResp.Data.ShortURL)
	return createResp.Data.Slug
}

func createCustomShortURL() {
	fmt.Println("=== Create Custom Short URL ===")
	// Use a unique slug to avoid conflicts in repeated runs
	customSlug := fmt.Sprintf("custom-%d", time.Now().Unix())
	expireAt := time.Now().Add(30 * 24 * time.Hour).Unix()

	customResp, err := client.CreateShortURL(seesdk.CreateShortURLRequest{
		TargetURL:  "https://www.example.com/custom",
		Domain:     defaultDomain,
		CustomSlug: customSlug,
		ExpireAt:   expireAt,
		Title:      "Custom Link",
		// TagIDs:     []int64{1, 2}, // Optional
	})
	if err != nil {
		log.Printf("Failed to create custom short URL: %v\n", err)
		return
	}
	fmt.Printf("Custom Slug: %s\n", customResp.Data.CustomSlug)
	fmt.Printf("Short URL: %s\n\n", customResp.Data.ShortURL)
}

func createProtectedShortURL() {
	fmt.Println("=== Create Password-Protected Short URL ===")
	protectedResp, err := client.CreateShortURL(seesdk.CreateShortURLRequest{
		TargetURL: "https://www.example.com/protected",
		Domain:    defaultDomain,
		Password:  "secret123",
		Title:     "Protected Link",
	})
	if err != nil {
		log.Printf("Failed to create protected short URL: %v\n", err)
		return
	}
	fmt.Printf("Protected Short URL: %s\n\n", protectedResp.Data.ShortURL)
}

func updateShortURL(slug string) {
	fmt.Println("=== Update Short URL ===")
	updateResp, err := client.UpdateShortURL(seesdk.UpdateShortURLRequest{
		Domain:    defaultDomain,
		Slug:      slug,
		TargetURL: "https://www.example.com/updated",
		Title:     "Updated Link",
	})
	if err != nil {
		log.Printf("Failed to update short URL: %v\n", err)
		return
	}
	fmt.Printf("Update successful: %s\n\n", updateResp.Message)
}

func deleteShortURL(slug string) {
	fmt.Println("=== Delete Short URL ===")
	deleteResp, err := client.DeleteShortURL(seesdk.DeleteURLRequest{
		Domain: defaultDomain,
		Slug:   slug,
	})
	if err != nil {
		log.Printf("Failed to delete URL: %v\n", err)
		return
	}
	fmt.Printf("Delete successful: %s\n", deleteResp.Message)
}

func fileOperations() {
	fmt.Println("=== File Operations ===")

	// Get file domains
	fileDomains, err := client.GetFileDomains()
	if err != nil {
		log.Printf("Failed to get file domains: %v\n", err)
	} else {
		fmt.Printf("File domains: %v\n", fileDomains.Data.Domains)
	}

	// Create a dummy file for upload
	tmpFile, err := os.CreateTemp("", "example-upload-*.txt")
	if err != nil {
		log.Printf("Failed to create temp file: %v\n", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	content := []byte(fmt.Sprintf("Hello S.EE SDK File Upload Test %d", time.Now().UnixNano()))
	if _, err := tmpFile.Write(content); err != nil {
		log.Printf("Failed to write to temp file: %v\n", err)
		return
	}
	tmpFile.Close()

	// Open file for reading
	fileToUpload, err := os.Open(tmpFile.Name())
	if err != nil {
		log.Printf("Failed to open file: %v\n", err)
		return
	}
	defer fileToUpload.Close()

	// Upload file
	fmt.Println("Uploading file...")
	uploadResp, err := client.UploadFile("test-file.txt", fileToUpload)
	if err != nil {
		log.Printf("Failed to upload file: %v\n", err)
		return
	}

	fmt.Printf("File uploaded successfully!\n")
	fmt.Printf("File URL: %s\n", uploadResp.Data.URL)
	fmt.Printf("Delete Key: %s\n", uploadResp.Data.Delete)

	// Delete file
	fmt.Println("Deleting file...")
	deleteFileResp, err := client.DeleteFile(uploadResp.Data.Hash)
	if err != nil {
		log.Printf("Failed to delete file: %v\n", err)
		return
	}
	fmt.Printf("Delete success: %v\n", deleteFileResp.Success)
}
