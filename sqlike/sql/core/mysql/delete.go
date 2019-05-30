package mysql

import (
	"github.com/si3nloong/sqlike/sqlike/actions"
	sqlstmt "github.com/si3nloong/sqlike/sqlike/sql/stmt"
)

// Delete :
func (ms *MySQL) Delete(f *actions.DeleteActions) (stmt *sqlstmt.Statement, err error) {
	stmt = sqlstmt.NewStatement(ms)
	err = buildStatement(stmt, ms.parser, f)
	if err != nil {
		return
	}
	return
}