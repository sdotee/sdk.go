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
// Last Modified: 2025-12-04 17:59:18
//

package seesdk

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

type CreateTextResponse struct {
	Code int `json:"code"`
	Data struct {
		CustomSlug string `json:"custom_slug"`
		ShortURL   string `json:"short_url"`
		Slug       string `json:"slug"`
	}
	Message string `json:"message"`
}

// GetUsageResponse represents the response containing usage statistics.
type GetUsageResponse struct {
	Code int `json:"code"`
	Data struct {
		APICountDay           int `json:"api_count_day"`
		APICountDayLimit      int `json:"api_count_day_limit"`
		APICountMonth         int `json:"api_count_month"`
		APICountMonthLimit    int `json:"api_count_month_limit"`
		LinkCountDay          int `json:"link_count_day"`
		LinkCountDayLimit     int `json:"link_count_day_limit"`
		LinkCountMonth        int `json:"link_count_month"`
		LinkCountMonthLimit   int `json:"link_count_month_limit"`
		QRCodeCountDay        int `json:"qrcode_count_day"`
		QRCodeCountDayLimit   int `json:"qrcode_count_day_limit"`
		QRCodeCountMonth      int `json:"qrcode_count_month"`
		QRCodeCountMonthLimit int `json:"qrcode_count_month_limit"`
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

// UploadFileResponse represents the response from uploading a file.
type UploadFileResponse struct {
	Code int `json:"code"`
	Data struct {
		Delete       string `json:"delete"`
		FileID       int    `json:"file_id"`
		Filename     string `json:"filename"`
		Hash         string `json:"hash"`
		Height       int    `json:"height"`
		Page         string `json:"page"`
		Path         string `json:"path"`
		Size         int    `json:"size"`
		Storename    string `json:"storename"`
		UploadStatus int    `json:"upload_status"`
		URL          string `json:"url"`
		Width        int    `json:"width"`
	} `json:"data"`
	Message string `json:"message"`
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
