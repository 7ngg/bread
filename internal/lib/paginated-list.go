package lib

type PaginatedList[T any] struct {
	Items      []T `json:"items"`
	Page       int32 `json:"page"`
	TotalPages int32 `json:"total_pages"`
}
