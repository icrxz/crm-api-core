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
