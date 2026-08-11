package request

import (
	"net/http"
	"strconv"
	"strings"
)

const (
	LimitPaginatedQuery  string = "limit"
	OffsetPaginatedQuery string = "offset"
	SortPaginatedQuery   string = "sort"
	SearchPaginatedQuery string = "search"
	TagsPaginatedQuery   string = "tags"
)

type CreatePostRequest struct {
	Title   string   `json:"title" validate:"required,max=200"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

type UpdatePostRequest struct {
	Title   *string  `json:"title" validate:"omitempty,min=1,max=500"`
	Content *string  `json:"content" validate:"omitempty,min=1,max=50000"`
	Tags    []string `json:"tags" validate:"max=20,dive,min=1,max=50"`
}
type PaginatedFeedQuery struct {
	Limit  int
	Offset int
	Sort   string
	Search string
	Tags   []string
}

func (q *PaginatedFeedQuery) Parse(r *http.Request) error {
	values := r.URL.Query()

	if value := values.Get(LimitPaginatedQuery); value != "" {
		limit, err := strconv.Atoi(value)
		if err != nil {
			return err
		}

		q.Limit = limit
	}

	if value := values.Get(OffsetPaginatedQuery); value != "" {
		offset, err := strconv.Atoi(value)
		if err != nil {
			return err
		}

		q.Offset = offset
	}

	if value := values.Get(SortPaginatedQuery); value != "" {
		q.Sort = strings.ToLower(value)
	}

	if value := values.Get(SearchPaginatedQuery); value != "" {
		q.Search = value
	}

	if value := values.Get(TagsPaginatedQuery); value != "" {
		q.Tags = strings.Split(value, ",")
	}

	return nil
}
