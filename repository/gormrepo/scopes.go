package gormrepo

import (
	"fmt"
	"strings"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ClausesV2(conds ...clause.Expression) func(tx *gorm.Statement) {
	return func(tx *gorm.Statement) {
		tx.Scopes(Clauses(conds...))
	}
}

// Clauses converts raw GORM clauses into a scope.
func Clauses(conds ...clause.Expression) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Clauses(conds...)
	}
}

// Offset adds an offset scope when offset is not nil.
func Offset(offset *int) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if offset != nil {
			return tx.Offset(*offset)
		}
		return tx
	}
}

// Limit adds a limit scope when limit is not nil.
func Limit(limit *int) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if limit != nil {
			return tx.Limit(*limit)
		}
		return tx
	}
}

// OrderASC adds an ascending order scope for the given columns.
func OrderASC(columns ...string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if len(columns) > 0 {
			return tx.Order(strings.Join(columns, ",") + " ASC")
		}
		return tx
	}
}

// OrderDESC adds a descending order scope for the given columns.
func OrderDESC(columns ...string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if len(columns) > 0 {
			return tx.Order(strings.Join(columns, ",") + " DESC")
		}
		return tx
	}
}

// Order adds a raw order scope for the given columns.
func Order(columns ...string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if len(columns) > 0 {
			return tx.Order(strings.Join(columns, ",") + " ")
		}
		return tx
	}
}

// Greater adds a column > value condition.
func Greater(column string, v any) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where(fmt.Sprintf("%v > ?", column), v)
	}
}

// GreaterEq adds a column >= value condition.
func GreaterEq(column string, v any) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where(fmt.Sprintf("%v >= ?", column), v)
	}
}

// Gt is a short alias for Greater.
func Gt(column string, v any) func(tx *gorm.DB) *gorm.DB {
	return Greater(column, v)
}

// Gte is a short alias for GreaterEq.
func Gte(column string, v any) func(tx *gorm.DB) *gorm.DB {
	return GreaterEq(column, v)
}

// Less adds a column < value condition.
func Less(column string, v any) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where(fmt.Sprintf("%v < ?", column), v)
	}
}

// LessEq adds a column <= value condition.
func LessEq(column string, v any) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where(fmt.Sprintf("%v <= ?", column), v)
	}
}

// Lt is an alias for Less.
func Lt(column string, v any) func(tx *gorm.DB) *gorm.DB {
	return Less(column, v)
}

// Lte is an alias for LessEq.
func Lte(column string, v any) func(tx *gorm.DB) *gorm.DB {
	return LessEq(column, v)
}

// Equal adds a column = value condition. With MySQL collations this is
// case-insensitive for string columns; use EqualStrict for binary comparison.
func Equal(column string, value any) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if len(column) > 0 && value != nil {
			return tx.Where(fmt.Sprintf("%s = ?", column), value)
		}
		return tx
	}
}

// EqualStrict adds a binary column = value condition.
func EqualStrict(column string, value any) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if len(column) > 0 && value != nil {
			return tx.Where(fmt.Sprintf("%s = binary ?", column), value)
		}
		return tx
	}
}

// NotEqual adds a column <> value condition.
func NotEqual(column string, value any) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if len(column) > 0 && value != nil {
			return tx.Where(fmt.Sprintf("%s <> ?", column), value)
		}
		return tx
	}
}

// IsNull adds a column IS NULL condition.
func IsNull(column string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if len(column) > 0 {
			return tx.Where(fmt.Sprintf("%s IS NULL", column))
		}
		return tx
	}
}

// IsNotNull adds a column IS NOT NULL condition.
func IsNotNull(column string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if len(column) > 0 {
			return tx.Where(fmt.Sprintf("%s IS NOT NULL", column))
		}
		return tx
	}
}

// In adds a column IN (values...) condition.
func In[T comparable](column string, values ...T) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if len(column) > 0 && len(values) > 0 {
			return tx.Where(fmt.Sprintf("%s in (?)", column), values)
		}
		return tx
	}
}

// NotIn adds a column NOT IN (values...) condition.
func NotIn[T comparable](column string, values ...T) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if len(column) > 0 && len(values) > 0 {
			return tx.Where(fmt.Sprintf("%s not in (?)", column), values)
		}
		return tx
	}
}

// Like adds a column LIKE %value% condition.
func Like(column string, value string) func(tx *gorm.DB) *gorm.DB {
	return LikePattern(column, fmt.Sprintf("%%%v%%", value))
}

// LikeEscaped adds a column LIKE %value% condition with LIKE wildcards escaped.
func LikeEscaped(column string, value string) func(tx *gorm.DB) *gorm.DB {
	return LikePatternEscaped(column, fmt.Sprintf("%%%v%%", escapeLikeLiteral(value)))
}

// LikePrefix adds a column LIKE value% condition.
func LikePrefix(column string, value string) func(tx *gorm.DB) *gorm.DB {
	return LikePattern(column, fmt.Sprintf("%v%%", value))
}

// LikeSuffix adds a column LIKE %value condition.
func LikeSuffix(column string, value string) func(tx *gorm.DB) *gorm.DB {
	return LikePattern(column, fmt.Sprintf("%%%v", value))
}

// LikePattern adds a raw LIKE pattern condition.
func LikePattern(column string, pattern string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if len(column) > 0 && pattern != "" {
			return tx.Where(fmt.Sprintf("%s like ?", column), pattern)
		}
		return tx
	}
}

// LikePatternEscaped adds a LIKE pattern condition with escape syntax.
func LikePatternEscaped(column string, pattern string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if len(column) > 0 && pattern != "" {
			return tx.Where(fmt.Sprintf("%s like ? escape '\\\\'", column), pattern)
		}
		return tx
	}
}

func escapeLikeLiteral(value string) string {
	replacer := strings.NewReplacer(
		"\\", "\\\\",
		"%", "\\%",
		"_", "\\_",
	)
	return replacer.Replace(value)
}

// Select adds a column selection scope.
func Select(columns ...string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if len(columns) > 0 {
			return tx.Select(columns)
		}
		return tx
	}
}

// Omit adds a column omission scope.
func Omit(columns ...string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if len(columns) > 0 {
			return tx.Omit(columns...)
		}
		return tx
	}
}

// Or combines conditions with OR.
func Or(conds ...func(tx *gorm.DB) *gorm.DB) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		oldtx := tx.Session(&gorm.Session{NewDB: true})
		for _, cond := range conds {
			oldtx = oldtx.Or(tx.Session(&gorm.Session{NewDB: true}).Scopes(cond))
		}
		return tx.Where(oldtx)
	}
}

// And combines conditions with AND.
func And(conds ...func(tx *gorm.DB) *gorm.DB) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tmptx := tx.Session(&gorm.Session{NewDB: true})
		for _, cond := range conds {
			tmptx = tmptx.Where(tx.Session(&gorm.Session{NewDB: true}).Scopes(cond))
		}
		return tx.Where(tmptx)
	}
}

// JsonColEqual adds a JSON column field equality condition.
func JsonColEqual(column string, k, v string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where(datatypes.JSONQuery(column).Equals(v, k))
	}
}
