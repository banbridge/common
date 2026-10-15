package gormrepo

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"

	"github.com/banbridge/common/repository"
)

// DefaultListLimit is applied when a list request only sets Offset. PGSQL and
// MariaDB require a LIMIT before OFFSET, so an explicit default is needed.
const DefaultListLimit = 10

// ConvertListOptions converts neutral ListOptions into GORM scopes.
func ConvertListOptions(opts repository.ListOptions) []func(tx *gorm.DB) *gorm.DB {
	var scopes []func(tx *gorm.DB) *gorm.DB

	if opts.Offset != nil {
		scopes = append(scopes, Offset(opts.Offset))
	}
	if opts.Limit != nil {
		scopes = append(scopes, Limit(opts.Limit))
	} else if opts.Offset != nil {
		// MySQL/MariaDB require a LIMIT before OFFSET; keep the original
		// default when only an offset was requested.
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

// ConvertConditions converts neutral conditions into GORM scopes. Conditions
// with an empty field or an unsupported operator are skipped.
func ConvertConditions(conds []repository.Condition) []func(tx *gorm.DB) *gorm.DB {
	var scopes []func(tx *gorm.DB) *gorm.DB
	for _, cond := range conds {
		if scope, ok := convertCondition(cond); ok {
			scopes = append(scopes, scope)
		}
	}
	return scopes
}

// ConvertConditionsStrict converts neutral conditions into GORM scopes and
// fails on empty fields or unsupported operators instead of silently dropping
// them. A dropped condition can turn an exact query into a full-table scan.
func ConvertConditionsStrict(conds []repository.Condition) ([]func(tx *gorm.DB) *gorm.DB, error) {
	scopes := make([]func(tx *gorm.DB) *gorm.DB, 0, len(conds))
	for i, cond := range conds {
		if cond.Field == "" {
			return nil, fmt.Errorf("gormrepo: condition field is empty at index %d", i)
		}
		scope, ok := convertCondition(cond)
		if !ok {
			return nil, fmt.Errorf("gormrepo: unsupported condition operator %q at index %d", cond.Operator, i)
		}
		scopes = append(scopes, scope)
	}
	return scopes, nil
}

func convertCondition(cond repository.Condition) (func(tx *gorm.DB) *gorm.DB, bool) {
	if cond.Field == "" {
		return nil, false
	}
	switch cond.Operator {
	case repository.Eq:
		return Equal(cond.Field, cond.Value), true
	case repository.In:
		return In(cond.Field, convertToAnySlice(cond.Value)...), true
	case repository.Gt:
		return Greater(cond.Field, cond.Value), true
	case repository.Gte:
		return GreaterEq(cond.Field, cond.Value), true
	case repository.Lt:
		return Less(cond.Field, cond.Value), true
	case repository.Lte:
		return LessEq(cond.Field, cond.Value), true
	case repository.NotEq:
		return NotEqual(cond.Field, cond.Value), true
	case repository.NotIn:
		return NotIn(cond.Field, convertToAnySlice(cond.Value)...), true
	default:
		return nil, false
	}
}

func convertToAnySlice(input any) []any {
	if input == nil {
		return []any{}
	}

	val := reflect.ValueOf(input)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		result := make([]any, val.Len())
		for i := 0; i < val.Len(); i++ {
			result[i] = val.Index(i).Interface()
		}
		return result
	}

	return []any{input}
}
