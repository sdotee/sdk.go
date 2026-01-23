//
// Copyright (c) 2025-2026 S.EE Development Team
//
// This source code is licensed under the MIT License,
// which is located in the LICENSE file in the source tree's root directory.
//
// File: api.go
// Author: S.EE Development Team <dev@s.ee>
// File Created: 2025-11-28 11:26:19
//
// Modified By: S.EE Development Team <dev@s.ee>
// Last Modified: 2025-12-04 17:58:55
//

package seesdk

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// unmarshalResponse is a helper function to unmarshal API response.
func unmarshalResponse[T any](data []byte, response *T) error {
	if err := json.Unmarshal(data, response); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}
	return nil
}

// CreateShortURL creates a new short URL with the given parameters.
func (c *Client) CreateShortURL(req CreateShortURLRequest) (*CreateShortURLResponse, error) {
	respBody, err := c.doRequest("POST", "/shorten", req)
	if err != nil {
		return nil, err
	}

	var response CreateShortURLResponse
	if err := unmarshalResponse(respBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// UpdateShortURL updates an existing short URL.
func (c *Client) UpdateShortURL(request UpdateShortURLRequest) (*UpdateShortURLResponse, error) {
	respBody, err := c.doRequest("PUT", "/shorten", request)
	if err != nil {
		return nil, err
	}

	var response UpdateShortURLResponse
	if err := unmarshalResponse(respBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteShortURL deletes an existing short URL.
func (c *Client) DeleteShortURL(request DeleteURLRequest) (*DeleteURLResponse, error) {
	respBody, err := c.doRequest("DELETE", "/shorten", request)
	if err != nil {
		return nil, err
	}

	var response DeleteURLResponse
	if err := unmarshalResponse(respBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetDomains retrieves the list of available domains.
func (c *Client) GetDomains() (*DomainsResponse, error) {
	respBody, err := c.doRequest("GET", "/domains", nil)
	if err != nil {
		return nil, err
	}

	var response DomainsResponse
	if err := unmarshalResponse(respBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetTags retrieves the list of available tags.
func (c *Client) GetTags() (*TagsResponse, error) {
	respBody, err := c.doRequest("GET", "/tags", nil)
	if err != nil {
		return nil, err
	}

	var response TagsResponse
	if err := unmarshalResponse(respBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateText creates a new text entry with the given parameters.
func (c *Client) CreateText(req CreateTextRequest) (*CreateTextResponse, error) {
	respBody, err := c.doRequest("POST", "/text", req)
	if err != nil {
		return nil, err
	}

	var response CreateTextResponse
	if err := unmarshalResponse(respBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// UpdateText updates an existing text entry.
func (c *Client) UpdateText(req UpdateTextRequest) (*UpdateTextResponse, error) {
	respBody, err := c.doRequest("PUT", "/text", req)
	if err != nil {
		return nil, err
	}

	var response UpdateTextResponse
	if err := unmarshalResponse(respBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteText deletes an existing text entry.
func (c *Client) DeleteText(req DeleteTextRequest) (*DeleteTextResponse, error) {
	respBody, err := c.doRequest("DELETE", "/text", req)
	if err != nil {
		return nil, err
	}

	var response DeleteTextResponse
	if err := unmarshalResponse(respBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// UploadFile uploads a file to the server.
func (c *Client) UploadFile(filename string, file io.Reader) (*UploadFileResponse, error) {
	if file == nil {
		return nil, fmt.Errorf("file is nil")
	}

	const maxFileSize = 100 * 1024 * 1024 // 100MB
	if err := checkFileSize(file, maxFileSize); err != nil {
		return nil, err
	}

	respBody, err := c.doMultipartRequest("/file/upload", "file", filename, file)
	if err != nil {
		return nil, err
	}

	var response UploadFileResponse
	if err := unmarshalResponse(respBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteFile deletes an uploaded file using its delete key.
func (c *Client) DeleteFile(deleteKey string) (*DeleteFileResponse, error) {
	respBody, err := c.doRequest("GET", "/file/delete/"+deleteKey, nil)
	if err != nil {
		return nil, err
	}

	var response DeleteFileResponse
	if err := unmarshalResponse(respBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetFileDomains retrieves the list of available domains for file sharing.
func (c *Client) GetFileDomains() (*DomainsResponse, error) {
	respBody, err := c.doRequest("GET", "/file/domains", nil)
	if err != nil {
		return nil, err
	}

	var response DomainsResponse
	if err := unmarshalResponse(respBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetTextDomains retrieves the list of available domains for text sharing.
func (c *Client) GetTextDomains() (*DomainsResponse, error) {
	respBody, err := c.doRequest("GET", "/text/domains", nil)
	if err != nil {
		return nil, err
	}

	var response DomainsResponse
	if err := unmarshalResponse(respBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// checkFileSize checks if the file size exceeds the maximum allowed size.
func checkFileSize(file io.Reader, maxSize int64) error {
	if f, ok := file.(interface{ Stat() (os.FileInfo, error) }); ok {
		if info, err := f.Stat(); err == nil && info.Size() > maxSize {
			return fmt.Errorf("file size exceeds the limit of %d bytes", maxSize)
		}
	} else if l, ok := file.(interface{ Len() int }); ok {
		if int64(l.Len()) > maxSize {
			return fmt.Errorf("file size exceeds the limit of %d bytes", maxSize)
		}
	}
	return nil
}
