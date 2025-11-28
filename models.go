package seesdk

// CreateShortURLRequest 创建短网址请求
type CreateShortURLRequest struct {
	CustomSlug            string  `json:"custom_slug,omitempty"`
	Domain                string  `json:"domain"`
	ExpirationRedirectURL string  `json:"expiration_redirect_url,omitempty"`
	ExpireAt              int64   `json:"expire_at,omitempty"` // unix timestamp (seconds)
	Password              string  `json:"password,omitempty"`
	TagIDs                []int64 `json:"tag_ids,omitempty"`
	TargetURL             string  `json:"target_url"`
	Title                 string  `json:"title,omitempty"`
}

// CreateShortURLResponse 创建短网址响应
type CreateShortURLResponse struct {
	Code int `json:"code"`
	Data struct {
		CustomSlug string `json:"custom_slug"`
		ShortURL   string `json:"short_url"`
		Slug       string `json:"slug"`
	} `json:"data"`
	Message string `json:"message"`
}

// DeleteURLRequest 删除短网址请求
type DeleteURLRequest struct {
	Domain string `json:"domain"`
	Slug   string `json:"slug"`
}

// DeleteURLResponse 删除短网址响应
type DeleteURLResponse struct {
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message"`
}

type UpdateShortURLRequest struct {
	Domain    string `json:"domain"`
	Slug      string `json:"slug"`
	TargetUrl string `json:"target_url"`
	Title     string `json:"title"`
}

type UpdateShortURLResponse struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type DomainsResponse struct {
	Code int `json:"code"`
	Data struct {
		Domains []string `json:"domains"`
	} `json:"data"`
	Message string `json:"message"`
}

type Tag struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type TagsResponse struct {
	Code int `json:"code"`
	Data struct {
		Tags []Tag `json:"tags"`
	} `json:"data"`
	Message string `json:"message"`
}
