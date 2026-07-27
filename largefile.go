//
// Copyright (c) 2025-2026 S.EE Development Team
//
// This source code is licensed under the MIT License,
// which is located in the LICENSE file in the source tree's root directory.
//
// File: largefile.go
// Author: S.EE Development Team <dev@s.ee>
// File Created: 2026-07-27 12:00:00
//
// Modified By: S.EE Development Team <dev@s.ee>
// Last Modified: 2026-07-27 12:00:00
//

package seesdk

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// maxLargeFileSize is the maximum allowed size for large file uploads (5GB).
const maxLargeFileSize = 5 * 1024 * 1024 * 1024

// DefaultLargeFileChunkSize is the chunk size used by UploadLargeFile (16MB).
const DefaultLargeFileChunkSize = 16 * 1024 * 1024

// tusProtocolVersion is the TUS protocol version supported by the S.EE API.
const tusProtocolVersion = "1.0.0"

// Large file upload session status values reported by GetLargeFileUploadProgress.
const (
	LargeFileUploadStatusUploading = 1
	LargeFileUploadStatusCompleted = 2
	LargeFileUploadStatusFailed    = 3
	LargeFileUploadStatusCancelled = 4
)

// CreateLargeFileUpload creates a TUS upload session for files up to 5GB.
// If FileHash is provided and the file already exists on the server, the
// response has FastUpload set to true and ExistingFile populated, and no
// data transfer is needed.
func (c *Client) CreateLargeFileUpload(req CreateLargeFileUploadRequest) (*CreateLargeFileUploadResponse, error) {
	if req.FileSize <= 0 {
		return nil, fmt.Errorf("file size must be positive")
	}
	if req.FileSize > maxLargeFileSize {
		return nil, fmt.Errorf("file size exceeds the limit of %d bytes", int64(maxLargeFileSize))
	}
	return callAPI[CreateLargeFileUploadResponse](c, "POST", "/file/large-file/create", req)
}

// CompleteLargeFileUpload finalizes an upload session after all chunks have
// been uploaded. It validates the file, moves it to permanent storage, and
// returns the file record. This consumes the upload session.
func (c *Client) CompleteLargeFileUpload(uploadID string) (*CompleteLargeFileUploadResponse, error) {
	return callAPI[CompleteLargeFileUploadResponse](c, "POST", "/file/large-file/complete", CompleteLargeFileUploadRequest{UploadID: uploadID})
}

// CancelLargeFileUpload cancels an in-progress upload session and removes
// the temporary data.
func (c *Client) CancelLargeFileUpload(uploadID string) (*CancelLargeFileUploadResponse, error) {
	return callAPI[CancelLargeFileUploadResponse](c, "DELETE", "/file/large-file/cancel", CancelLargeFileUploadRequest{UploadID: uploadID})
}

// GetLargeFileUploadProgress returns the current progress of a large file upload session.
func (c *Client) GetLargeFileUploadProgress(uploadID string) (*GetLargeFileUploadProgressResponse, error) {
	return callAPI[GetLargeFileUploadProgressResponse](c, "GET", "/file/large-file/progress?upload_id="+url.QueryEscape(uploadID), nil)
}

// GetLargeFileUploadOffset queries the server for the current upload offset
// of a session (TUS HEAD request). Use it to resume an interrupted upload.
func (c *Client) GetLargeFileUploadOffset(uploadID string) (int64, error) {
	req, err := http.NewRequest(http.MethodHead, c.tusEndpoint(uploadID), nil)
	if err != nil {
		return 0, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.doTUS(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	offset, err := strconv.ParseInt(resp.Header.Get("Upload-Offset"), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse Upload-Offset header: %w", err)
	}
	return offset, nil
}

// UploadLargeFileChunk uploads a single chunk at the given offset (TUS PATCH
// request) and returns the new offset reported by the server.
func (c *Client) UploadLargeFileChunk(uploadID string, offset int64, chunk []byte) (int64, error) {
	req, err := http.NewRequest(http.MethodPatch, c.tusEndpoint(uploadID), bytes.NewReader(chunk))
	if err != nil {
		return 0, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/offset+octet-stream")
	req.Header.Set("Upload-Offset", strconv.FormatInt(offset, 10))

	resp, err := c.doTUS(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	newOffset, err := strconv.ParseInt(resp.Header.Get("Upload-Offset"), 10, 64)
	if err != nil {
		// Server did not report an offset; assume the whole chunk was accepted.
		return offset + int64(len(chunk)), nil
	}
	return newOffset, nil
}

// UploadLargeFile uploads a file up to 5GB using the TUS resumable protocol.
// It creates an upload session, transfers the content in chunks of
// DefaultLargeFileChunkSize, and completes the session. When the server
// reports an instant (deduplicated) upload, no data is transferred and the
// existing file record is returned.
func (c *Client) UploadLargeFile(req CreateLargeFileUploadRequest, r io.Reader) (*CompleteLargeFileUploadResponse, error) {
	if r == nil {
		return nil, fmt.Errorf("file is nil")
	}

	createResp, err := c.CreateLargeFileUpload(req)
	if err != nil {
		return nil, err
	}
	if createResp.Code != 200 {
		return nil, fmt.Errorf("create upload session failed (code %d): %s", createResp.Code, createResp.Message)
	}

	// Instant upload: the file already exists on the server.
	if createResp.Data.FastUpload && createResp.Data.ExistingFile != nil {
		return &CompleteLargeFileUploadResponse{
			Code:    createResp.Code,
			Data:    CompleteLargeFileUploadData{File: *createResp.Data.ExistingFile},
			Message: createResp.Message,
		}, nil
	}

	uploadID := createResp.Data.UploadID
	buf := make([]byte, DefaultLargeFileChunkSize)
	var offset int64
	for offset < req.FileSize {
		n, err := io.ReadFull(r, buf)
		if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
			_, _ = c.CancelLargeFileUpload(uploadID)
			return nil, fmt.Errorf("read file content: %w", err)
		}

		offset, err = c.UploadLargeFileChunk(uploadID, offset, buf[:n])
		if err != nil {
			_, _ = c.CancelLargeFileUpload(uploadID)
			return nil, err
		}
	}

	return c.CompleteLargeFileUpload(uploadID)
}

// tusEndpoint returns the TUS endpoint URL for an upload session.
func (c *Client) tusEndpoint(uploadID string) string {
	return c.BaseURL + "/file/large-file-tus/" + url.PathEscape(uploadID)
}

// doTUS sets TUS and common headers, executes the request, and validates the
// response status. The caller must close the response body on success.
func (c *Client) doTUS(req *http.Request) (*http.Response, error) {
	req.Header.Set("Tus-Resumable", tusProtocolVersion)
	req.Header.Set("User-Agent", "see-go-sdk/"+Version)
	if c.APIKey != "" {
		req.Header.Set("Authorization", c.APIKey)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("TUS error (status %d): %s", resp.StatusCode, string(body))
	}
	return resp, nil
}
