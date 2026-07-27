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
// Last Modified: 2026-07-27 12:00:00
//

package seesdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
)

// UsageNoLimit represents unlimited usage.
const UsageNoLimit = -1

// maxUploadFileSize is the maximum allowed size for file uploads (100MB).
const maxUploadFileSize = 100 * 1024 * 1024

// Visit statistics period values accepted by GetLinkVisitStat.
const (
	VisitStatPeriodDaily   = "daily"   // today
	VisitStatPeriodMonthly = "monthly" // this month
	VisitStatPeriodTotally = "totally" // all-time (default)
)

// callAPI executes a JSON API request and unmarshals the response into T.
func callAPI[T any](c *Client, method, endpoint string, body any) (*T, error) {
	respBody, err := c.doRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}
	return unmarshalResponse[T](respBody)
}

// unmarshalResponse unmarshals an API response body into T.
func unmarshalResponse[T any](data []byte) (*T, error) {
	var response T
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &response, nil
}

// pagedEndpoint appends the page query parameter when page is greater than 1.
func pagedEndpoint(endpoint string, page int) string {
	if page > 1 {
		return fmt.Sprintf("%s?page=%d", endpoint, page)
	}
	return endpoint
}

// CreateShortURL creates a new short URL with the given parameters.
func (c *Client) CreateShortURL(req CreateShortURLRequest) (*CreateShortURLResponse, error) {
	return callAPI[CreateShortURLResponse](c, "POST", "/shorten", req)
}

// UpdateShortURL updates an existing short URL.
func (c *Client) UpdateShortURL(req UpdateShortURLRequest) (*UpdateShortURLResponse, error) {
	return callAPI[UpdateShortURLResponse](c, "PUT", "/shorten", req)
}

// DeleteShortURL deletes an existing short URL.
func (c *Client) DeleteShortURL(req DeleteURLRequest) (*DeleteURLResponse, error) {
	return callAPI[DeleteURLResponse](c, "DELETE", "/shorten", req)
}

// GetLinkVisitStat retrieves click/visit statistics for a short URL.
// The period can be VisitStatPeriodDaily, VisitStatPeriodMonthly, or
// VisitStatPeriodTotally. An empty period defaults to all-time statistics.
func (c *Client) GetLinkVisitStat(domain, slug, period string) (*GetLinkVisitStatResponse, error) {
	query := url.Values{}
	query.Set("domain", domain)
	query.Set("slug", slug)
	if period != "" {
		query.Set("period", period)
	}
	return callAPI[GetLinkVisitStatResponse](c, "GET", "/link/visit-stat?"+query.Encode(), nil)
}

// GetUsage retrieves the usage statistics of the account.
func (c *Client) GetUsage() (*GetUsageResponse, error) {
	return callAPI[GetUsageResponse](c, "GET", "/usage", nil)
}

// GetDomains retrieves the list of available domains.
func (c *Client) GetDomains() (*DomainsResponse, error) {
	return callAPI[DomainsResponse](c, "GET", "/domains", nil)
}

// GetTags retrieves the list of available tags.
func (c *Client) GetTags() (*TagsResponse, error) {
	return callAPI[TagsResponse](c, "GET", "/tags", nil)
}

// CreateText creates a new text entry with the given parameters.
func (c *Client) CreateText(req CreateTextRequest) (*CreateTextResponse, error) {
	return callAPI[CreateTextResponse](c, "POST", "/text", req)
}

// UpdateText updates an existing text entry.
func (c *Client) UpdateText(req UpdateTextRequest) (*UpdateTextResponse, error) {
	return callAPI[UpdateTextResponse](c, "PUT", "/text", req)
}

// DeleteText deletes an existing text entry.
func (c *Client) DeleteText(req DeleteTextRequest) (*DeleteTextResponse, error) {
	return callAPI[DeleteTextResponse](c, "DELETE", "/text", req)
}

// UploadFile uploads a file to the server.
func (c *Client) UploadFile(req UploadFileRequest) (*UploadFileResponse, error) {
	if req.File == nil {
		return nil, fmt.Errorf("file is nil")
	}

	if err := checkFileSize(req.File, maxUploadFileSize); err != nil {
		return nil, err
	}

	fields := make(map[string]string)
	if req.Domain != "" {
		fields["domain"] = req.Domain
	}
	if req.CustomSlug != "" {
		fields["custom_slug"] = req.CustomSlug
	}
	if req.IsPrivate {
		fields["is_private"] = "1"
	}

	respBody, err := c.doMultipartRequest("/file/upload", "file", req.Filename, req.File, fields)
	if err != nil {
		return nil, err
	}
	return unmarshalResponse[UploadFileResponse](respBody)
}

// SmartUploadFile uploads a file choosing the best strategy automatically:
// files up to 100MB go through the regular multipart upload, while larger
// files (up to 5GB) are transferred with the TUS resumable protocol. The
// reader should expose its size via Stat() (e.g. *os.File) or Len() (e.g.
// *bytes.Reader); when the size cannot be determined, the regular upload
// is used.
func (c *Client) SmartUploadFile(req UploadFileRequest) (*UploadFileResponse, error) {
	if req.File == nil {
		return nil, fmt.Errorf("file is nil")
	}

	size, ok := readerSize(req.File)
	if !ok || size <= maxUploadFileSize {
		return c.UploadFile(req)
	}

	isPrivate := 0
	if req.IsPrivate {
		isPrivate = 1
	}

	largeResp, err := c.UploadLargeFile(CreateLargeFileUploadRequest{
		Alias:     req.CustomSlug,
		Domain:    req.Domain,
		FileName:  req.Filename,
		FileSize:  size,
		IsPrivate: isPrivate,
	}, req.File)
	if err != nil {
		return nil, err
	}

	return &UploadFileResponse{
		Code:    largeResp.Code,
		Data:    largeResp.Data.File,
		Message: largeResp.Message,
	}, nil
}

// GetFileHistory retrieves a paginated list of uploaded files.
// Returns 30 files per page, sorted by creation time descending.
// Page starts at 1. If page is 0 or negative, defaults to page 1.
func (c *Client) GetFileHistory(page int) (*GetFileHistoryResponse, error) {
	return callAPI[GetFileHistoryResponse](c, "GET", pagedEndpoint("/files", page), nil)
}

// GetLinkHistory retrieves a paginated list of short links created by the account.
// Page starts at 1. If page is 0 or negative, defaults to page 1.
func (c *Client) GetLinkHistory(page int) (*GetLinkHistoryResponse, error) {
	return callAPI[GetLinkHistoryResponse](c, "GET", pagedEndpoint("/links", page), nil)
}

// GetTextHistory retrieves a paginated list of text sharings created by the account.
// Page starts at 1. If page is 0 or negative, defaults to page 1.
func (c *Client) GetTextHistory(page int) (*GetTextHistoryResponse, error) {
	return callAPI[GetTextHistoryResponse](c, "GET", pagedEndpoint("/texts", page), nil)
}

// CheckToken checks whether the provided API token is valid and usable.
// When valid, the response includes the token string and its expiration time.
func (c *Client) CheckToken(token string) (*CheckTokenResponse, error) {
	return callAPI[CheckTokenResponse](c, "POST", "/token/check", CheckTokenRequest{Token: token})
}

// DeleteFile deletes an uploaded file using its delete key.
func (c *Client) DeleteFile(deleteKey string) (*DeleteFileResponse, error) {
	return callAPI[DeleteFileResponse](c, "GET", "/file/delete/"+url.PathEscape(deleteKey), nil)
}

// GetFileDomains retrieves the list of available domains for file sharing.
func (c *Client) GetFileDomains() (*DomainsResponse, error) {
	return callAPI[DomainsResponse](c, "GET", "/file/domains", nil)
}

// GetTextDomains retrieves the list of available domains for text sharing.
func (c *Client) GetTextDomains() (*DomainsResponse, error) {
	return callAPI[DomainsResponse](c, "GET", "/text/domains", nil)
}

// GetPrivateFileDownloadURL retrieves a temporary download URL for a private file using its file ID.
// The URL is valid for a limited time (about 1 hour).
func (c *Client) GetPrivateFileDownloadURL(fileID int64) (*GetPrivateFileDownloadURLResponse, error) {
	endpoint := fmt.Sprintf("/file/private/download-url?file_id=%d", fileID)
	return callAPI[GetPrivateFileDownloadURLResponse](c, "GET", endpoint, nil)
}

// readerSize reports the size of r when it exposes Stat() (e.g. *os.File)
// or Len() (e.g. *bytes.Reader, *bytes.Buffer, *strings.Reader).
func readerSize(r io.Reader) (int64, bool) {
	if f, ok := r.(interface{ Stat() (os.FileInfo, error) }); ok {
		if info, err := f.Stat(); err == nil {
			return info.Size(), true
		}
		return 0, false
	}
	if l, ok := r.(interface{ Len() int }); ok {
		return int64(l.Len()), true
	}
	return 0, false
}

// checkFileSize checks if the file size exceeds the maximum allowed size.
func checkFileSize(file io.Reader, maxSize int64) error {
	if size, ok := readerSize(file); ok && size > maxSize {
		return fmt.Errorf("file size exceeds the limit of %d bytes", maxSize)
	}
	return nil
}
