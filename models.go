package seesdk

// CreateShortURLRequest represents a request to create a short URL
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

// CreateShortURLResponse represents the response from creating a short URL
type CreateShortURLResponse struct {
	Code int `json:"code"`
	Data struct {
		CustomSlug string `json:"custom_slug"`
		ShortURL   string `json:"short_url"`
		Slug       string `json:"slug"`
	} `json:"data"`
	Message string `json:"message"`
}

// DeleteURLRequest represents a request to delete a short URL
type DeleteURLRequest struct {
	Domain string `json:"domain"`
	Slug   string `json:"slug"`
}

// DeleteURLResponse represents the response from deleting a short URL
type DeleteURLResponse struct {
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message"`
}

// UpdateShortURLRequest represents a request to update a short URL
type UpdateShortURLRequest struct {
	Domain    string `json:"domain"`
	Slug      string `json:"slug"`
	TargetURL string `json:"target_url"`
	Title     string `json:"title"`
}

// UpdateShortURLResponse represents the response from updating a short URL
type UpdateShortURLResponse struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

// DomainsResponse represents the response containing available domains
type DomainsResponse struct {
	Code int `json:"code"`
	Data struct {
		Domains []string `json:"domains"`
	} `json:"data"`
	Message string `json:"message"`
}

// Tag represents a tag entity
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// TagsResponse represents the response containing available tags
type TagsResponse struct {
	Code int `json:"code"`
	Data struct {
		Tags []Tag `json:"tags"`
	} `json:"data"`
	Message string `json:"message"`
}
