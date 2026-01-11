package domain

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

type PagingFilter struct {
	Limit  int
	Offset int
	SortBy string
	SortOrder SortOrder
}

type Paging struct {
	Total  int
	Limit  int
	Offset int
}

type PagingResult[T any] struct {
	Result []T
	Paging Paging
}
