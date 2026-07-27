//
// Copyright (c) 2025-2026 S.EE Development Team
//
// This source code is licensed under the MIT License,
// which is located in the LICENSE file in the source tree's root directory.
//
// File: qrcode.go
// Author: S.EE Development Team <dev@s.ee>
// File Created: 2026-07-27 12:00:00
//
// Modified By: S.EE Development Team <dev@s.ee>
// Last Modified: 2026-07-27 12:00:00
//

package seesdk

// CreateQRCode creates a dynamic QR code for a given URL. The QR code is
// generated server-side and available in PNG, SVG, and PDF formats.
func (c *Client) CreateQRCode(req CreateQRCodeRequest) (*CreateQRCodeResponse, error) {
	return callAPI[CreateQRCodeResponse](c, "POST", "/qrcode", req)
}

// DeleteQRCode permanently deletes a QR code and its associated short link.
func (c *Client) DeleteQRCode(req DeleteQRCodeRequest) (*DeleteQRCodeResponse, error) {
	return callAPI[DeleteQRCodeResponse](c, "DELETE", "/qrcode", req)
}

// GetQRCodeHistory retrieves a paginated list of QR codes.
// Page starts at 1. If page is 0 or negative, defaults to page 1.
func (c *Client) GetQRCodeHistory(page int) (*GetQRCodeHistoryResponse, error) {
	return callAPI[GetQRCodeHistoryResponse](c, "GET", pagedEndpoint("/qrcodes", page), nil)
}
