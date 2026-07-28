//
// Copyright (c) 2025 S.EE Development Team
//
// This source code is licensed under the MIT License,
// which is located in the LICENSE file in the source tree's root directory.
//
// File: models.go
// Author: S.EE Development Team <dev@s.ee>
// File Created: 2025-11-28 11:26:17
//
// Modified By: S.EE Development Team <dev@s.ee>
// Last Modified: 2026-07-27 12:00:00
//

package seesdk

import "io"

// CreateShortURLRequest represents a request to create a short URL.
type CreateShortURLRequest struct {
	CustomSlug            string  `json:"custom_slug,omitempty"`
	Domain                string  `json:"domain"`
	ExpirationRedirectURL string  `json:"expiration_redirect_url,omitempty"`
	ExpireAt              int64   `json:"expire_at,omitempty"` // Unix timestamp in seconds
	Password              string  `json:"password,omitempty"`
	TagIDs                []int64 `json:"tag_ids,omitempty"`
	TargetURL             string  `json:"target_url"`
	Title                 string  `json:"title,omitempty"`
}

type CreateTextRequest struct {
	Content    string  `json:"content"`
	CustomSlug string  `json:"custom_slug,omitempty"`
	Domain     string  `json:"domain,omitempty"`
	ExpireAt   int64   `json:"expire_at,omitempty"` // Unix timestamp in seconds
	Password   string  `json:"password,omitempty"`
	TagIDs     []int64 `json:"tag_ids,omitempty"`
	TextType   string  `json:"text_type,omitempty"`
	Title      string  `json:"title,omitempty"`
}

// CreateShortURLResponse represents the response from creating a short URL.
type CreateShortURLResponse struct {
	Code int `json:"code"`
	Data struct {
		CustomSlug string `json:"custom_slug"`
		ShortURL   string `json:"short_url"`
		Slug       string `json:"slug"`
	} `json:"data"`
	Message string `json:"message"`
}

// CreateTextResponse represents the response from creating a text sharing.
type CreateTextResponse struct {
	Code int `json:"code"`
	Data struct {
		CustomSlug string `json:"custom_slug"`
		ShortURL   string `json:"short_url"`
		Slug       string `json:"slug"`
	} `json:"data"`
	Message string `json:"message"`
}

// GetUsageResponse represents the response containing usage statistics.
type GetUsageResponse struct {
	Code int `json:"code"`
	Data struct {
		APICountDay           int    `json:"api_count_day"`
		APICountDayLimit      int    `json:"api_count_day_limit"`
		APICountMonth         int    `json:"api_count_month"`
		APICountMonthLimit    int    `json:"api_count_month_limit"`
		LinkCountDay          int    `json:"link_count_day"`
		LinkCountDayLimit     int    `json:"link_count_day_limit"`
		LinkCountMonth        int    `json:"link_count_month"`
		LinkCountMonthLimit   int    `json:"link_count_month_limit"`
		QRCodeCountDay        int    `json:"qrcode_count_day"`
		QRCodeCountDayLimit   int    `json:"qrcode_count_day_limit"`
		QRCodeCountMonth      int    `json:"qrcode_count_month"`
		QRCodeCountMonthLimit int    `json:"qrcode_count_month_limit"`
		TextCountDay          int    `json:"text_count_day"`
		TextCountDayLimit     int    `json:"text_count_day_limit"`
		TextCountMonth        int    `json:"text_count_month"`
		TextCountMonthLimit   int    `json:"text_count_month_limit"`
		UploadCountDay        int    `json:"upload_count_day"`
		UploadCountDayLimit   int    `json:"upload_count_day_limit"`
		UploadCountMonth      int    `json:"upload_count_month"`
		UploadCountMonthLimit int    `json:"upload_count_month_limit"`
		FileCount             int    `json:"file_count"`
		StorageUsageMB        string `json:"storage_usage_mb"`       // in MB, rounded to 2 decimal places
		StorageUsageLimitMB   string `json:"storage_usage_limit_mb"` // in MB, "-1" means unlimited
	} `json:"data"`
	Message string `json:"message"`
}

// GetLinkVisitStatResponse represents the response containing link visit statistics.
type GetLinkVisitStatResponse struct {
	Code int `json:"code"`
	Data struct {
		VisitCount int64 `json:"visit_count"`
	} `json:"data"`
	Message string `json:"message"`
}

// DeleteURLRequest represents a request to delete a short URL.
type DeleteURLRequest struct {
	Domain string `json:"domain"`
	Slug   string `json:"slug"`
}

// DeleteTextRequest represents a request to delete a text.
type DeleteTextRequest struct {
	Domain string `json:"domain"`
	Slug   string `json:"slug"`
}

// DeleteURLResponse represents the response from deleting a short URL.
type DeleteURLResponse struct {
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message"`
}

// DeleteTextResponse represents the response from deleting a text.
type DeleteTextResponse struct {
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message"`
}

// UpdateShortURLRequest represents a request to update a short URL.
type UpdateShortURLRequest struct {
	Domain    string `json:"domain"`
	Slug      string `json:"slug"`
	TargetURL string `json:"target_url"`
	Title     string `json:"title"`
}

// UploadFileRequest represents a request to upload a file.
type UploadFileRequest struct {
	Filename   string
	File       io.Reader
	Domain     string
	CustomSlug string
	IsPrivate  bool
}

// UploadFileData represents the metadata of an uploaded file.
type UploadFileData struct {
	CreatedAt    int    `json:"created_at,omitempty"`
	Delete       string `json:"delete"`
	FileID       int    `json:"file_id"`
	Filename     string `json:"filename"`
	Hash         string `json:"hash"`
	Height       int    `json:"height"`
	MimeType     string `json:"mime_type,omitempty"`
	Page         string `json:"page"`
	Path         string `json:"path"`
	Size         int    `json:"size"`
	Storename    string `json:"storename"`
	ThumbURL     string `json:"thumb_url,omitempty"`
	UploadStatus int    `json:"upload_status"`
	URL          string `json:"url"`
	Width        int    `json:"width"`
}

// UploadFileResponse represents the response from uploading a file.
type UploadFileResponse struct {
	Code    int            `json:"code"`
	Data    UploadFileData `json:"data"`
	Message string         `json:"message"`
}

// GetFileHistoryResponse represents the response containing file upload history.
type GetFileHistoryResponse struct {
	Code    int              `json:"code"`
	Data    []UploadFileData `json:"data"`
	Message string           `json:"message"`
	Success bool             `json:"success"`
}

// DeleteFileResponse represents the response from deleting a file.
type DeleteFileResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type UpdateTextRequest struct {
	Domain  string `json:"domain"`
	Slug    string `json:"slug"`
	Content string `json:"content"`
	Title   string `json:"title,omitempty"`
}

// UpdateShortURLResponse represents the response from updating a short URL.
type UpdateShortURLResponse struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

// UpdateTextResponse represents the response from updating a text.
type UpdateTextResponse struct {
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message"`
}

// DomainsResponse represents the response containing available domains.
type DomainsResponse struct {
	Code int `json:"code"`
	Data struct {
		Domains []string `json:"domains"`
	} `json:"data"`
	Message string `json:"message"`
}

// Tag represents a tag entity.
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// TagsResponse represents the response containing available tags.
type TagsResponse struct {
	Code int `json:"code"`
	Data struct {
		Tags []Tag `json:"tags"`
	} `json:"data"`
	Message string `json:"message"`
}

// GetPrivateFileDownloadURLData contains the private file download URL information.
type GetPrivateFileDownloadURLData struct {
	FileID    int64  `json:"file_id"`
	URL       string `json:"url"`
	ExpiresAt int64  `json:"expires_at"`
}

// GetPrivateFileDownloadURLResponse represents the response when getting a private file download URL.
type GetPrivateFileDownloadURLResponse struct {
	Code    int                           `json:"code"`
	Data    GetPrivateFileDownloadURLData `json:"data"`
	Message string                        `json:"message"`
	Success bool                          `json:"success"`
}

// CreateLargeFileUploadRequest represents a request to create a large file (TUS) upload session.
type CreateLargeFileUploadRequest struct {
	Alias       string `json:"alias,omitempty"` // Custom slug for the short link (alphanumeric only)
	Description string `json:"description,omitempty"`
	Domain      string `json:"domain,omitempty"`
	ExpireAt    int64  `json:"expire_at,omitempty"` // Unix timestamp for link expiry
	FileHash    string `json:"file_hash,omitempty"` // SHA256 hash for instant upload deduplication
	FileName    string `json:"file_name"`
	FileSize    int64  `json:"file_size"`
	IsPrivate   int    `json:"is_private,omitempty"` // 0 = public (default), 1 = private
	MimeType    string `json:"mime_type,omitempty"`
	Password    string `json:"password,omitempty"`
	Title       string `json:"title,omitempty"`
}

// CreateLargeFileUploadData contains the created upload session information.
type CreateLargeFileUploadData struct {
	ExistingFile *UploadFileData `json:"existing_file,omitempty"` // Populated when FastUpload is true
	ExpiresAt    int64           `json:"expires_at"`              // Unix timestamp when this upload session expires (24h)
	FastUpload   bool            `json:"fast_upload"`             // True when the file already exists and upload was skipped
	FileSize     int64           `json:"file_size"`
	ID           int64           `json:"id"`
	UploadID     string          `json:"upload_id"`
	UploadURL    string          `json:"upload_url"` // TUS upload endpoint URL
}

// CreateLargeFileUploadResponse represents the response from creating a large file upload session.
type CreateLargeFileUploadResponse struct {
	Code    int                       `json:"code"`
	Data    CreateLargeFileUploadData `json:"data"`
	Message string                    `json:"message"`
}

// CompleteLargeFileUploadRequest represents a request to complete a large file upload session.
type CompleteLargeFileUploadRequest struct {
	UploadID string `json:"upload_id"`
}

// CompleteLargeFileUploadData contains the completed upload result.
type CompleteLargeFileUploadData struct {
	File      UploadFileData `json:"file"`
	ShortLink string         `json:"short_link,omitempty"` // Short link created for the file (if domain/alias were provided)
}

// CompleteLargeFileUploadResponse represents the response from completing a large file upload.
type CompleteLargeFileUploadResponse struct {
	Code    int                         `json:"code"`
	Data    CompleteLargeFileUploadData `json:"data"`
	Message string                      `json:"message"`
}

// CancelLargeFileUploadRequest represents a request to cancel a large file upload session.
type CancelLargeFileUploadRequest struct {
	UploadID string `json:"upload_id"`
}

// CancelLargeFileUploadResponse represents the response from cancelling a large file upload.
type CancelLargeFileUploadResponse struct {
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message"`
}

// LargeFileUploadProgressData contains the progress of a large file upload session.
type LargeFileUploadProgressData struct {
	CreatedAt    int64   `json:"created_at"`
	FileName     string  `json:"file_name"`
	FileSize     int64   `json:"file_size"`
	Progress     float64 `json:"progress"` // Progress percentage (0-100)
	Status       int     `json:"status"`   // See LargeFileUploadStatus* constants
	UpdatedAt    int64   `json:"updated_at"`
	UploadID     string  `json:"upload_id"`
	UploadedSize int64   `json:"uploaded_size"`
}

// GetLargeFileUploadProgressResponse represents the response containing upload progress.
type GetLargeFileUploadProgressResponse struct {
	Code    int                         `json:"code"`
	Data    LargeFileUploadProgressData `json:"data"`
	Message string                      `json:"message"`
}

// LinkHistoryData represents a short link entry in the creation history.
type LinkHistoryData struct {
	CreatedAt  int64  `json:"created_at"`
	Domain     string `json:"domain"`
	ObjectType int    `json:"object_type"` // 0 = link, 1 = file, 2 = text, 3 = qrcode
	ShortURL   string `json:"short_url"`
	Slug       string `json:"slug"`
	TargetURL  string `json:"target_url"`
	Title      string `json:"title"`
	VisitCount int64  `json:"visit_count"`
}

// GetLinkHistoryResponse represents the response containing link creation history.
type GetLinkHistoryResponse struct {
	Code    int               `json:"code"`
	Data    []LinkHistoryData `json:"data"`
	Message string            `json:"message"`
	Success bool              `json:"success"`
}

// TextHistoryData represents a text sharing entry in the creation history.
type TextHistoryData struct {
	ContentPreview string `json:"content_preview"` // Preview of the content (truncated)
	CreatedAt      int64  `json:"created_at"`
	Domain         string `json:"domain"`
	ID             int64  `json:"id"`
	IsExpired      bool   `json:"is_expired"`
	ShortURL       string `json:"short_url"`
	Slug           string `json:"slug"`
	TextType       string `json:"text_type"` // "plain_text", "source_code", or "markdown"
	Title          string `json:"title"`
}

// GetTextHistoryResponse represents the response containing text creation history.
type GetTextHistoryResponse struct {
	Code    int               `json:"code"`
	Data    []TextHistoryData `json:"data"`
	Message string            `json:"message"`
	Success bool              `json:"success"`
}

// CheckTokenRequest represents a request to validate an API token.
type CheckTokenRequest struct {
	Token string `json:"token"`
}

// CheckTokenResponse represents the response from validating an API token.
type CheckTokenResponse struct {
	Code int `json:"code"`
	Data struct {
		ExpiresAt int64  `json:"expires_at"` // Token expiration time in Unix timestamp
		Token     string `json:"token"`
		Valid     bool   `json:"valid"`
	} `json:"data"`
	Message string `json:"message"`
}

// BioCustomLink represents a custom link on a bio page.
type BioCustomLink struct {
	Description string `json:"description,omitempty"`
	Title       string `json:"title"`
	URL         string `json:"url"`
}

// CreateBioPageRequest represents a request to create a bio page.
type CreateBioPageRequest struct {
	CustomLinks []BioCustomLink `json:"custom_links,omitempty"`
	CustomSlug  string          `json:"custom_slug,omitempty"`
	Description string          `json:"description,omitempty"`
	Domain      string          `json:"domain,omitempty"`
	MastodonURL string          `json:"mastodon_url,omitempty"`
	RSSURL      string          `json:"rss_url,omitempty"`
	Title       string          `json:"title"`
}

// CreateBioPageResponse represents the response from creating a bio page.
type CreateBioPageResponse struct {
	Code int `json:"code"`
	Data struct {
		BioPageID int64  `json:"bio_page_id"`
		ShortURL  string `json:"short_url"`
	} `json:"data"`
	Message string `json:"message"`
}

// UpdateBioPageRequest represents a request to update a bio page.
// Omit CustomLinks to leave the existing custom links unchanged.
type UpdateBioPageRequest struct {
	CustomLinks []BioCustomLink `json:"custom_links,omitempty"`
	Description string          `json:"description,omitempty"`
	ID          int64           `json:"id"`
	MastodonURL string          `json:"mastodon_url,omitempty"`
	RSSURL      string          `json:"rss_url,omitempty"`
	Title       string          `json:"title"`
}

// UpdateBioPageResponse represents the response from updating a bio page.
type UpdateBioPageResponse struct {
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message"`
}

// DeleteBioPageRequest represents a request to delete a bio page.
type DeleteBioPageRequest struct {
	ID int64 `json:"id"`
}

// DeleteBioPageResponse represents the response from deleting a bio page.
type DeleteBioPageResponse struct {
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message"`
}

// BioPageData represents a bio page entry in the history.
type BioPageData struct {
	CreatedAt   int64           `json:"created_at"`
	CustomLinks []BioCustomLink `json:"custom_links"`
	Description string          `json:"description"`
	Domain      string          `json:"domain"`
	ID          int64           `json:"id"`
	Link        string          `json:"link"` // The complete short URL of the bio page
	MastodonURL string          `json:"mastodon_url"`
	RSSURL      string          `json:"rss_url"`
	Slug        string          `json:"slug"`
	Title       string          `json:"title"`
}

// GetBioPageHistoryResponse represents the response containing bio page history.
type GetBioPageHistoryResponse struct {
	Code int `json:"code"`
	Data struct {
		BioPages []BioPageData `json:"bio_pages"`
		Total    int64         `json:"total"`
	} `json:"data"`
	Message string `json:"message"`
}

// CreateQRCodeRequest represents a request to create a dynamic QR code.
type CreateQRCodeRequest struct {
	CustomSlug string `json:"custom_slug,omitempty"`
	Domain     string `json:"domain,omitempty"`
	TargetURL  string `json:"target_url"`
	Title      string `json:"title"`
}

// CreateQRCodeData contains the created QR code information.
type CreateQRCodeData struct {
	CustomSlug string `json:"custom_slug"`
	PDFURL     string `json:"pdf_url"`
	PNGURL     string `json:"png_url"`
	ShortURL   string `json:"short_url"`
	Slug       string `json:"slug"`
	SVGURL     string `json:"svg_url"`
}

// CreateQRCodeResponse represents the response from creating a QR code.
type CreateQRCodeResponse struct {
	Code    int              `json:"code"`
	Data    CreateQRCodeData `json:"data"`
	Message string           `json:"message"`
}

// DeleteQRCodeRequest represents a request to delete a QR code.
type DeleteQRCodeRequest struct {
	Domain string `json:"domain"`
	Slug   string `json:"slug"`
}

// DeleteQRCodeResponse represents the response from deleting a QR code.
type DeleteQRCodeResponse struct {
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message"`
}

// QRCodeHistoryData represents a QR code entry in the history.
type QRCodeHistoryData struct {
	CreatedAt int64  `json:"created_at"`
	Domain    string `json:"domain"`
	PDFURL    string `json:"pdf_url"`
	PNGURL    string `json:"png_url"`
	ScanCount int64  `json:"scan_count"`
	ShortURL  string `json:"short_url"`
	Slug      string `json:"slug"`
	SVGURL    string `json:"svg_url"`
	Title     string `json:"title"`
}

// GetQRCodeHistoryResponse represents the response containing QR code history.
type GetQRCodeHistoryResponse struct {
	Code int `json:"code"`
	Data struct {
		QRCodes []QRCodeHistoryData `json:"qrcodes"`
		Total   int64               `json:"total"`
	} `json:"data"`
	Message string `json:"message"`
}
