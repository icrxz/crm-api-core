package database

import (
	"fmt"
	"strings"
)

func prepareInQuery[S comparable](filters []S, query []string, args []any, key string) ([]string, []any) {
	if len(filters) > 0 {
		parsedArray := make([]any, 0, len(filters))
		for _, filter := range filters {
			parsedArray = append(parsedArray, filter)
		}

		queryFormatted := fmt.Sprintf("%s IN (", key)
		for i := len(args) + 1; i < len(args)+1+len(filters); i++ {
			queryFormatted += fmt.Sprintf("$%d,", i)
		}
		queryFormatted = strings.TrimRight(queryFormatted, ",")
		queryFormatted += ")"

		query = append(query, queryFormatted)
		args = append(args, parsedArray...)
	}

	return query, args
}

func prepareLikeQuery[S comparable](filters []S, query []string, args []any, key string) ([]string, []any) {
	if len(filters) > 0 {
		parsedArray := make([]any, 0, len(filters))
		for _, filter := range filters {
			parsedArray = append(parsedArray, filter)
		}

		queryFormatted := key + " LIKE '%' || "
		for i := len(args) + 1; i < len(args)+1+len(filters); i++ {
			queryFormatted += fmt.Sprintf("$%d::text", i)
		}
		queryFormatted += " || '%'"

		query = append(query, queryFormatted)
		args = append(args, parsedArray...)
	}

	return query, args
}

func createChunks[T any](slice []T, size int) [][]T {
	var chunks [][]T
	for i := 0; i < len(slice); i += size {
		end := i + size

		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}
