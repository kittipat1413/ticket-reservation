package entity

import "fmt"

var (
	ErrInvalidSortOrder = fmt.Errorf("invalid sort order")
)

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

var sortOrderStringMapper = map[SortOrder]string{
	SortOrderAsc:  "asc",
	SortOrderDesc: "desc",
}

func (s SortOrder) String() string {
	return sortOrderStringMapper[s]
}

func (s SortOrder) IsValid() bool {
	switch s {
	case SortOrderAsc, SortOrderDesc:
		return true
	default:
		return false
	}
}

// Parse parses a string into a SortOrder. It returns an error if the string is not a valid SortOrder.
func ParseSortOrder(order string) (SortOrder, error) {
	sortOrder := SortOrder(order)
	if !sortOrder.IsValid() {
		return "", fmt.Errorf("%w: %s", ErrInvalidSortOrder, order)
	}
	return sortOrder, nil
}
