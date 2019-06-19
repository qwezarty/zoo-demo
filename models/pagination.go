package models

type Pagination struct {
	PageSize  int         `json:"page_size"`
	PageToken int         `json:"page_token"`
	TotalSize int         `json:"total_size"`
	Body      interface{} `json:"body"`
}
