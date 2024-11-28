package debug

import (
	"github.com/Oskang09/sqlike/sql/dialect"
	"github.com/Oskang09/sqlike/sql/dialect/mysql"
	sqlstmt "github.com/Oskang09/sqlike/sql/stmt"
)

// ToSQL :
func ToSQL(src interface{}) error {
	ms := dialect.GetDialectByDriver("mysql").(*mysql.MySQL)
	sqlstmt.NewStatement(ms)
	return nil
}
