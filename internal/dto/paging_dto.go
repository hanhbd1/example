package dto

type PageRequest struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

type PagingDTO struct {
	Page       int         `json:"page"`
	Size       int         `json:"size"`
	Count      int         `json:"count"`
	TotalCount int         `json:"total_count"`
	Data       interface{} `json:"data"`
}
