package mysql

import (
	"github.com/si3nloong/sqlike/sqlike/actions"
	sqlstmt "github.com/si3nloong/sqlike/sqlike/sql/stmt"
)

// Update :
func (ms *MySQL) Update(f *actions.UpdateActions) (stmt *sqlstmt.Statement, err error) {
	stmt = sqlstmt.NewStatement(ms)
	err = buildStatement(stmt, ms.parser, f)
	if err != nil {
		return
	}
	return
}