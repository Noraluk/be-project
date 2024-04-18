package models

import "github.com/gofiber/fiber/v2"

type Response struct {
	Data interface{} `json:"data"`
}

func (m Response) ToPagination(c *fiber.Ctx, totalRecords int64) PaginationResponse {
	totalPages := totalRecords / int64(c.QueryInt("page_size"))
	if totalRecords%int64(c.QueryInt("page_size")) != 0 {
		totalPages += 1
	}
	page := c.QueryInt("page")
	if page == 0 {
		page = 1
	}
	return PaginationResponse{
		Response:     m,
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		CurrentPage:  page,
	}
}

type PaginationResponse struct {
	Response
	TotalRecords int64 `json:"total_records"`
	CurrentPage  int   `json:"current_page"`
	TotalPages   int64 `json:"total_pages"`
}
