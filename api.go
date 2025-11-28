package seesdk

import (
	"encoding/json"
	"fmt"
)

// CreateShortURL 创建短网址
func (c *Client) CreateShortURL(req CreateShortURLRequest) (*CreateShortURLResponse, error) {
	respBody, err := c.doRequest("POST", "/shorten", req)
	if err != nil {
		return nil, err
	}

	var response CreateShortURLResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &response, nil
}

func (c *Client) UpdateShortURL(request UpdateShortURLRequest) (*UpdateShortURLResponse, error) {
	respBody, err := c.doRequest("PUT", "/shorten", request)
	if err != nil {
		return nil, err
	}

	var response UpdateShortURLResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &response, nil
}

// DeleteURL 删除短网址
func (c *Client) DeleteShortURL(request DeleteURLRequest) (*DeleteURLResponse, error) {
	respBody, err := c.doRequest("DELETE", "/shorten", request)
	if err != nil {
		return nil, err
	}

	var response DeleteURLResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &response, nil
}

func (c *Client) GetDomains() (*DomainsResponse, error) {
	respBody, err := c.doRequest("GET", "/domains", nil)
	if err != nil {
		return nil, err
	}

	var response DomainsResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &response, nil
}

func (c *Client) GetTags() (*TagsResponse, error) {
	respBody, err := c.doRequest("GET", "/tags", nil)
	if err != nil {
		return nil, err
	}

	var response TagsResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &response, nil
}
