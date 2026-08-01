package types

type CloudflareListRecordsResult struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Content    string                 `json:"content"`
	Proxiable  bool                   `json:"proxiable"`
	Proxied    bool                   `json:"proxied"`
	TTL        int                    `json:"ttl"`
	Settings   map[string]interface{} `json:"settings"`
	Meta       map[string]interface{} `json:"meta"`
	Comment    string                 `json:"comment"`
	Tags       []string               `json:"tags"`
	CreatedOn  string                 `json:"created_on"`
	ModifiedOn string                 `json:"modified_on"`
}

type CloudflareErrors struct {
	Code             int    `json:"code"`
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
	Source           struct {
		Pointer string `json:"pointer"`
	} `json:"source"`
}

type CloudflareResultInfo struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Count      int `json:"count"`
	TotalCount int `json:"total_count"`
	TotalPages int `json:"total_pages"`
}

type CloudflareListRecordsResponse struct {
	Result     []CloudflareListRecordsResult `json:"result"`
	Success    bool                          `json:"success"`
	Errors     []CloudflareErrors            `json:"errors"`
	Messages   []string                      `json:"messages"`
	ResultInfo CloudflareResultInfo          `json:"result_info"`
}
