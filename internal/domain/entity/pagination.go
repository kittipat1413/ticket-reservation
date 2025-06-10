package entity

type PageProvider[T any] func() ([]T, PageProvider[T], Pagination, error)

type Page[T any] interface {
	GetPagination() Pagination
	GetData() []T
	Next() (Page[T], error)
}

type Pagination struct {
	Total       int64 `json:"total"`
	Limit       int64 `json:"limit"`
	Offset      int64 `json:"offset"`
	PageCount   int64 `json:"page_count,omitempty"`
	CurrentPage int64 `json:"current_page,omitempty"`
}

func NewPagination(total, limit, offset int64) Pagination {
	pageCount := total / limit
	if total%limit != 0 {
		pageCount++
	}
	currentPage := offset/limit + 1

	return Pagination{
		Total:       total,
		Limit:       limit,
		Offset:      offset,
		PageCount:   pageCount,
		CurrentPage: currentPage,
	}
}

type emptyPage[T any] struct{}

func (p emptyPage[T]) GetPagination() Pagination { return Pagination{} }
func (p emptyPage[T]) GetData() []T              { return []T{} }
func (p emptyPage[T]) Next() (Page[T], error)    { return emptyPage[T]{}, nil }

type pageImpl[T any] struct {
	provider   PageProvider[T]
	pagination Pagination
	data       []T
}

func (p *pageImpl[T]) GetPagination() Pagination { return p.pagination }
func (p *pageImpl[T]) GetData() []T              { return p.data }
func (p *pageImpl[T]) Next() (Page[T], error) {
	if data, nextProvider, pagination, err := p.provider(); err != nil {
		return nil, err
	} else if nextProvider != nil {
		return &pageImpl[T]{nextProvider, pagination, data}, nil
	} else {
		return emptyPage[T]{}, nil
	}
}

func NewPage[T any](prov PageProvider[T]) (Page[T], error) {
	if data, nextProvider, pagination, err := prov(); err != nil {
		return nil, err
	} else if len(data) > 0 {
		return &pageImpl[T]{nextProvider, pagination, data}, nil
	} else {
		return emptyPage[T]{}, nil
	}
}

func GetAllPageData[T any](page Page[T]) ([]T, error) {
	data := page.GetData()
	for len(page.GetData()) > 0 {
		if nextPage, err := page.Next(); err != nil {
			return nil, err
		} else {
			page = nextPage
		}
		data = append(data, page.GetData()...)
	}
	return data, nil
}
