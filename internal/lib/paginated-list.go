package lib

type PaginatedList[T any] struct {
	Items      []T `json:"items"`
	Page       int64 `json:"page"`
	TotalPages int64 `json:"total_pages"`
}
