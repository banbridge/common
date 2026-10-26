package gormrepo

import (
	"strings"

	"gorm.io/gorm"

	"github.com/banbridge/common/repository"
)

// DefaultListLimit is applied when a list request only sets Offset. PGSQL and
// MariaDB require a LIMIT before OFFSET, so an explicit default is needed.
const DefaultListLimit = 10

func ApplyListOptions(opts repository.ListOptions) func(db *gorm.Statement) {
	return func(db *gorm.Statement) {
		if opts.Offset != nil {
			db.Offset(*opts.Offset)
		}
		if opts.Limit != nil {
			db.Limit(*opts.Limit)

		} else if opts.Offset != nil {

			defaultLimit := DefaultListLimit

			db.Limit(defaultLimit)
		}

		for _, order := range opts.OrderBy {
			if order == nil {
				continue
			}
			t, cs := order()
			if len(cs) == 0 {
				continue
			}
			if t.IsASC() {
				db.Order(strings.Join(cs, ",") + " ASC")

			} else if t.IsDesc() {
				db.Order(strings.Join(cs, ",") + " DESC")
			}
		}
		if len(opts.SelectColumns) > 0 {
			db.Select(opts.SelectColumns)
		}

	}

}

func ApplyResetOffset() func(db *gorm.Statement) {
	return func(db *gorm.Statement) {
		db.Offset(-1).Limit(-1)
	}
}

// ConvertListOptions converts neutral ListOptions into GORM scopes.
func ConvertListOptions(opts repository.ListOptions) []func(tx *gorm.DB) *gorm.DB {
	var scopes []func(tx *gorm.DB) *gorm.DB

	if opts.Offset != nil {
		scopes = append(scopes, Offset(opts.Offset))
	}
	if opts.Limit != nil {
		scopes = append(scopes, Limit(opts.Limit))
	} else if opts.Offset != nil {
		defaultLimit := DefaultListLimit
		scopes = append(scopes, Limit(&defaultLimit))
	}

	for _, order := range opts.OrderBy {
		if order == nil {
			continue
		}
		t, cs := order()
		if len(cs) == 0 {
			continue
		}
		if t.IsASC() {
			scopes = append(scopes, OrderASC(cs...))
		} else if t.IsDesc() {
			scopes = append(scopes, OrderDESC(cs...))
		}
	}
	if len(opts.SelectColumns) > 0 {
		scopes = append(scopes, Select(opts.SelectColumns...))
	}
	return scopes
}
