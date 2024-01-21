package domain

import (
	"fmt"
	"strings"
)

type PaginatedResult struct {
	Data        any `json:"data"`
	CurrentPage int `json:"current_page"`
	Count       int `json:"count"`
	PageCount   int `json:"page_count"`
}

type StorageFilters struct {
	IsPaginated    bool
	IsSortedByDate bool

	Pagination PaginationFilters
	SortDate   SortDateFilters
}

func (filters *StorageFilters) BuildSQL(base string) string {
	buffer := strings.Builder{}
	buffer.WriteString(base)

	if filters.IsSortedByDate {
		buffer.WriteString(" ORDER BY created_at ")
		if filters.SortDate.IsASC {
			buffer.WriteString("ASC")
		} else {
			buffer.WriteString("DESC")
		}
	}

	if filters.IsPaginated {
		buffer.WriteString(fmt.Sprintf(
			" LIMIT %d OFFSET %d",
			filters.Pagination.Count,
			(filters.Pagination.Page-1)*filters.Pagination.Count,
		))
	}

	return buffer.String()
}

type PaginationFilters struct {
	Page  int
	Count int
}

type SortDateFilters struct {
	IsASC bool
}
