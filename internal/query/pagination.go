package query

import "github.com/Moaz125-eng/logforge/pkg/logentry"

type Page struct {
	Items      []logentry.Entry `json:"items"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}

func Paginate(entries []logentry.Entry, page, size int) Page {
	if size <= 0 {
		size = 50
	}
	if page <= 0 {
		page = 1
	}
	total := len(entries)
	totalPages := (total + size - 1) / size
	if totalPages == 0 {
		totalPages = 1
	}
	start := (page - 1) * size
	if start > total {
		start = total
	}
	end := start + size
	if end > total {
		end = total
	}
	return Page{
		Items:      entries[start:end],
		Total:      total,
		Page:       page,
		PageSize:   size,
		TotalPages: totalPages,
	}
}
