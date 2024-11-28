package mysql

import (
	sqlstmt "github.com/Oskang09/sqlike/sql/stmt"
	"github.com/Oskang09/sqlike/sqlike/actions"
)

// Delete :
func (ms *MySQL) Delete(stmt sqlstmt.Stmt, f *actions.DeleteActions) (err error) {
	err = buildStatement(stmt, ms.parser, f)
	if err != nil {
		return
	}
	return
}
