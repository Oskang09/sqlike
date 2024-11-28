package casbin

import (
	"github.com/Oskang09/sqlike/sql/expr"
	"github.com/Oskang09/sqlike/sqlike/primitive"
)

// Filter :
func Filter(fields ...interface{}) primitive.Group {
	return expr.And(fields...)
}
