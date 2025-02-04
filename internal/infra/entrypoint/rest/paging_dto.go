package rest

import "github.com/icrxz/crm-api-core/internal/domain"

type PagingFilterDTO struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type PagingDTO struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type SearchResultDTO[T any] struct {
	Result []T       `json:"result"`
	Paging PagingDTO `json:"paging"`
}

func mapPagingToPagingDTO(paging domain.Paging) PagingDTO {
	return PagingDTO{
		Total:  paging.Total,
		Limit:  paging.Limit,
		Offset: paging.Offset,
	}
}

func mapSearchResultToSearchResultDTO[T any, U any](result domain.PagingResult[T], mapFunc func(obj []T) []U) SearchResultDTO[U] {
	return SearchResultDTO[U]{
		Result: mapFunc(result.Result),
		Paging: mapPagingToPagingDTO(result.Paging),
	}
}
