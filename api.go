package seesdk

import (
	"encoding/json"
	"fmt"
)

// unmarshalResponse is a helper function to unmarshal API response
func unmarshalResponse[T any](data []byte, response *T) error {
	if err := json.Unmarshal(data, response); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}
	return nil
}

// CreateShortURL creates a new short URL with the given parameters
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

// UpdateShortURL updates an existing short URL
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

// DeleteShortURL deletes an existing short URL
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

// GetDomains retrieves the list of available domains
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

// GetTags retrieves the list of available tags
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
