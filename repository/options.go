// Package repository defines database-agnostic repository contracts used by
// business services and by GORM-based repository implementations.
package repository

// OrderType describes the sort direction of an OrderBy.
type OrderType string

const (
	// Asc sorts the selected columns in ascending order.
	Asc OrderType = "asc"
	// Desc sorts the selected columns in descending order.
	Desc OrderType = "desc"
)

// IsASC reports whether the order direction is ascending.
func (o OrderType) IsASC() bool {
	return o == Asc
}

// IsDesc reports whether the order direction is descending.
func (o OrderType) IsDesc() bool {
	return o == Desc
}

// OrderBy returns the sort direction and columns for a ListOptions query.
// Callers normally construct it with AscBy or DescBy.
type OrderBy func() (OrderType, []string)

// AscBy builds an ascending OrderBy for the given columns.
func AscBy(col ...string) OrderBy {
	return func() (OrderType, []string) {
		return Asc, col
	}
}

// DescBy builds a descending OrderBy for the given columns.
func DescBy(col ...string) OrderBy {
	return func() (OrderType, []string) {
		return Desc, col
	}
}

// ListOptions describes pagination, ordering and field selection for a list
// query. Nil pointers mean "not set"; Limit can be -1 to disable paging.
type ListOptions struct {
	Offset        *int
	Limit         *int
	OrderBy       []OrderBy
	SelectColumns []string
}
