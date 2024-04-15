package pagination

import (
	"fmt"
	"math"
	"net/url"
	"strings"

	"github.com/vekshinnikita/golang_vehicles/internal/tools"
	"github.com/vekshinnikita/golang_vehicles/internal/validators"
)

type Pagination struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Sort     string `json:"sort"`
}

type PaginatedResponse[T any] struct {
	FirstPage    int `json:"first_page"`
	LastPage     int `json:"last_page"`
	CurrentPage  int `json:"current_page"`
	TotalRecords int `json:"total_records"`
	PageSize     int `json:"page_size"`
	Data         T   `json:"data"`
}

func NewPagination(qs url.Values, v *validators.Validator) *Pagination {
	validatesPage := []validators.ValidateFunc[int]{validators.ValidateNumber(v, 0, 1)}

	page := tools.ReadIntUrlQuery(qs, v, validatesPage, "page", 1)
	pageSize := tools.ReadIntUrlQuery(qs, v, validatesPage, "page_size", 20)
	sort := tools.ReadStrUrlQuery(qs, "sort", "")

	return &Pagination{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
	}
}

func (p *Pagination) GetSQLSortQuery() string {
	order := "ASC"

	if p.Sort == "" {
		return "ORDER BY id ASC"
	}

	if strings.HasPrefix(p.Sort, "-") {
		order = "DESC"
	}

	field := strings.TrimPrefix(p.Sort, "-")

	return fmt.Sprintf("ORDER BY %s %s, id ASC", field, order)
}

func (p *Pagination) GetSQLPaginationQuery() string {
	offset := p.PageSize * (p.Page - 1)

	return fmt.Sprintf(`LIMIT %d OFFSET %d`, p.PageSize, offset)
}

func GetPaginatedResponse[T any](p Pagination, data T, totalRecords int) PaginatedResponse[T] {

	if totalRecords <= 0 {
		return PaginatedResponse[T]{}
	}

	return PaginatedResponse[T]{
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(p.PageSize))),
		CurrentPage:  p.Page,
		TotalRecords: totalRecords,
		PageSize:     p.PageSize,
		Data:         data,
	}
}
