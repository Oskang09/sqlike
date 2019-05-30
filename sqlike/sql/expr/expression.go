package expr

import (
	"github.com/si3nloong/sqlike/sqlike/primitive"
)

// Equal :
func Equal(field string, value interface{}) (c primitive.C) {
	c.Field = primitive.L(field)
	c.Operator = primitive.Equal
	c.Values = append(c.Values, value)
	return
}

// NotEqual :
func NotEqual(field string, value interface{}) (c primitive.C) {
	c.Field = primitive.L(field)
	c.Operator = primitive.NotEqual
	c.Values = append(c.Values, value)
	return
}

// IsNull :
func IsNull(field string) (c primitive.C) {
	c.Field = primitive.L(field)
	c.Operator = primitive.IsNull
	return
}

// NotNull :
func NotNull(field string) (c primitive.C) {
	c.Field = primitive.L(field)
	c.Operator = primitive.NotNull
	return
}

// In :
func In(field string, value interface{}) (c primitive.C) {
	c.Field = primitive.L(field)
	c.Operator = primitive.In
	c.Values = append(c.Values, value)
	return
}

// NotIn :
func NotIn(field string, value interface{}) (c primitive.C) {
	c.Field = primitive.L(field)
	c.Operator = primitive.NotIn
	c.Values = append(c.Values, value)
	return
}

// Like :
func Like(field string, value interface{}) (c primitive.C) {
	c.Field = primitive.L(field)
	c.Operator = primitive.Like
	c.Values = append(c.Values, value)
	return
}

// NotLike :
func NotLike(field string, value interface{}) (c primitive.C) {
	c.Field = primitive.L(field)
	c.Operator = primitive.NotLike
	c.Values = append(c.Values, value)
	return
}

// GreaterEqual :
func GreaterEqual(field string, value interface{}) (c primitive.C) {
	c.Field = primitive.L(field)
	c.Operator = primitive.GreaterEqual
	c.Values = append(c.Values, value)
	return
}

// GreaterThan :
func GreaterThan(field string, value interface{}) (c primitive.C) {
	c.Field = primitive.L(field)
	c.Operator = primitive.GreaterThan
	c.Values = append(c.Values, value)
	return
}

// LowerEqual :
func LowerEqual(field string, value interface{}) (c primitive.C) {
	c.Field = primitive.L(field)
	c.Operator = primitive.LowerEqual
	c.Values = append(c.Values, value)
	return
}

// LowerThan :
func LowerThan(field string, value interface{}) (c primitive.C) {
	c.Field = primitive.L(field)
	c.Operator = primitive.LowerThan
	c.Values = append(c.Values, value)
	return
}

// Between :
func Between(field string, from, to interface{}) (c primitive.C) {
	c.Field = primitive.L(field)
	c.Operator = primitive.Between
	c.Values = append(c.Values, from, to)
	return
}

// NotBetween :
func NotBetween(field string, from, to interface{}) (c primitive.C) {
	c.Field = primitive.L(field)
	c.Operator = primitive.NotBetween
	c.Values = append(c.Values, from, to)
	return
}

// And :
func And(conds ...interface{}) (g primitive.G) {
	if len(conds) > 1 {
		g = append(g, primitive.Raw(`(`))
		for i, cond := range conds {
			if i > 0 {
				g = append(g, primitive.And)
			}
			g = append(g, cond)
		}
		g = append(g, primitive.Raw(`)`))
		return
	}
	g = append(g, conds...)
	return
}

// Or :
func Or(conds ...interface{}) (g primitive.G) {
	if len(conds) > 1 {
		g = append(g, primitive.Raw(`(`))
		for i, cond := range conds {
			if i > 0 {
				g = append(g, primitive.Or)
			}
			g = append(g, primitive.Raw(`(`), cond, primitive.Raw(`)`))
		}
		g = append(g, primitive.Raw(`)`))
		return
	}
	g = append(g, conds...)
	return
}

// Asc :
func Asc(field string) (s primitive.Sort) {
	s.Field = field
	s.Order = primitive.Ascending
	return
}

// Desc :
func Desc(field string) (s primitive.Sort) {
	s.Field = field
	s.Order = primitive.Descending
	return
}