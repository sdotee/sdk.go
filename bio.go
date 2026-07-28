//
// Copyright (c) 2025-2026 S.EE Development Team
//
// This source code is licensed under the MIT License,
// which is located in the LICENSE file in the source tree's root directory.
//
// File: bio.go
// Author: S.EE Development Team <dev@s.ee>
// File Created: 2026-07-27 12:00:00
//
// Modified By: S.EE Development Team <dev@s.ee>
// Last Modified: 2026-07-27 12:00:00
//

package seesdk

// CreateBioPage creates a new bio page with a short URL.
func (c *Client) CreateBioPage(req CreateBioPageRequest) (*CreateBioPageResponse, error) {
	return callAPI[CreateBioPageResponse](c, "POST", "/bio", req)
}

// UpdateBioPage updates an existing bio page.
func (c *Client) UpdateBioPage(req UpdateBioPageRequest) (*UpdateBioPageResponse, error) {
	return callAPI[UpdateBioPageResponse](c, "PUT", "/bio", req)
}

// DeleteBioPage permanently deletes a bio page by its numeric ID.
func (c *Client) DeleteBioPage(id int64) (*DeleteBioPageResponse, error) {
	return callAPI[DeleteBioPageResponse](c, "DELETE", "/bio", DeleteBioPageRequest{ID: id})
}

// GetBioPageHistory retrieves a paginated list of bio pages.
// Page starts at 1. If page is 0 or negative, defaults to page 1.
func (c *Client) GetBioPageHistory(page int) (*GetBioPageHistoryResponse, error) {
	return callAPI[GetBioPageHistoryResponse](c, "GET", pagedEndpoint("/bios", page), nil)
}
