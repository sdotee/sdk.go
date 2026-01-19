//
// Copyright (c) 2025 S.EE Development Team
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

func main() {
	// Create SDK client
	client := seesdk.NewClient(seesdk.Config{
		BaseURL: "https://api.example.com", // Replace with actual API URL
		APIKey:  "your-api-key-here",       // Replace with actual API key
		Timeout: 30 * time.Second,
	})

	// Example 1: Get available domains
	fmt.Println("=== Get Domains ===")
	domainsResp, err := client.GetDomains()
	if err != nil {
		log.Fatalf("Failed to get domains: %v", err)
	}
	fmt.Printf("Available domains: %v\n\n", domainsResp.Data.Domains)

	// Example 2: Get available tags
	fmt.Println("=== Get Tags ===")
	tagsResp, err := client.GetTags()
	if err != nil {
		log.Fatalf("Failed to get tags: %v", err)
	}
	fmt.Printf("Available tags:\n")
	for _, tag := range tagsResp.Data.Tags {
		fmt.Printf("  - ID: %d, Name: %s\n", tag.ID, tag.Name)
	}
	fmt.Println()

	// Example 3: Create short URL with basic settings
	fmt.Println("=== Create Short URL ===")
	createResp, err := client.CreateShortURL(seesdk.CreateShortURLRequest{
		TargetURL: "https://www.example.com/very/long/url/path",
		Domain:    "example.com", // Use one of the available domains
		Title:     "Example Link",
	})
	if err != nil {
		log.Fatalf("Failed to create short URL: %v", err)
	}
	fmt.Printf("Response Code: %d\n", createResp.Code)
	fmt.Printf("Message: %s\n", createResp.Message)
	fmt.Printf("Slug: %s\n", createResp.Data.Slug)
	fmt.Printf("Short URL: %s\n\n", createResp.Data.ShortURL)

	// Example 4: Create custom short URL with expiration
	fmt.Println("=== Create Custom Short URL ===")
	expireAt := time.Now().Add(30 * 24 * time.Hour).Unix() // Expires in 30 days
	customResp, err := client.CreateShortURL(seesdk.CreateShortURLRequest{
		TargetURL:  "https://www.example.com/custom",
		Domain:     "example.com",
		CustomSlug: "my-custom-code",
		ExpireAt:   expireAt,
		Title:      "Custom Link",
		TagIDs:     []int64{1, 2}, // Use actual tag IDs
	})
	if err != nil {
		log.Fatalf("Failed to create custom short URL: %v", err)
	}
	fmt.Printf("Custom Slug: %s\n", customResp.Data.CustomSlug)
	fmt.Printf("Short URL: %s\n\n", customResp.Data.ShortURL)

	// Example 5: Create password-protected short URL
	fmt.Println("=== Create Password-Protected Short URL ===")
	protectedResp, err := client.CreateShortURL(seesdk.CreateShortURLRequest{
		TargetURL: "https://www.example.com/protected",
		Domain:    "example.com",
		Password:  "secret123",
		Title:     "Protected Link",
	})
	if err != nil {
		log.Fatalf("Failed to create protected short URL: %v", err)
	}
	fmt.Printf("Protected Short URL: %s\n\n", protectedResp.Data.ShortURL)

	// Example 6: Update short URL
	fmt.Println("=== Update Short URL ===")
	updateResp, err := client.UpdateShortURL(seesdk.UpdateShortURLRequest{
		Domain:    "example.com",
		Slug:      createResp.Data.Slug,
		TargetURL: "https://www.example.com/updated",
		Title:     "Updated Link",
	})
	if err != nil {
		log.Fatalf("Failed to update short URL: %v", err)
	}
	fmt.Printf("Update successful: %s\n\n", updateResp.Message)

	// Example 7: Delete short URL
	fmt.Println("=== Delete Short URL ===")
	deleteResp, err := client.DeleteShortURL(seesdk.DeleteURLRequest{
		Domain: "example.com",
		Slug:   createResp.Data.Slug,
	})
	if err != nil {
		log.Fatalf("Failed to delete URL: %v", err)
	}
	fmt.Printf("Delete successful: %s\n", deleteResp.Message)

	// Example 8: File Operations
	fmt.Println("\n=== File Operations ===")

	// Get file domains
	fileDomains, err := client.GetFileDomains()
	if err != nil {
		log.Printf("Failed to get file domains: %v", err)
	} else {
		fmt.Printf("File domains: %v\n", fileDomains.Data.Domains)
	}

	// Create a dummy file for upload
	tmpFile, err := os.CreateTemp("", "example-upload-*.txt")
	if err != nil {
		log.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := []byte(fmt.Sprintf("Hello S.EE SDK File Upload Test %d", time.Now().UnixNano()))
	if _, err := tmpFile.Write(content); err != nil {
		log.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close() // Close to reopen for reading

	// Open file for reading
	fileToUpload, err := os.Open(tmpFile.Name())
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer fileToUpload.Close()

	// Upload file
	fmt.Println("Uploading file...")
	uploadResp, err := client.UploadFile("test-file.txt", fileToUpload)
	if err != nil {
		log.Printf("Failed to upload file: %v", err)
	} else {
		fmt.Printf("File uploaded successfully!\n")
		fmt.Printf("File URL: %s\n", uploadResp.Data.URL)
		fmt.Printf("Delete Key: %s\n", uploadResp.Data.Delete)

		// Delete file
		fmt.Println("Deleting file...")
		deleteFileResp, err := client.DeleteFile(uploadResp.Data.Delete)
		if err != nil {
			log.Printf("Failed to delete file: %v", err)
		} else {
			fmt.Printf("Delete success: %v\n", deleteFileResp.Success)
		}
	}
}
