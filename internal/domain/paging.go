package domain

type PagingFilter struct {
	Limit  int
	Offset int
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
