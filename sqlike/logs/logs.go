package logs

import (
	sqlstmt "github.com/Oskang09/sqlike/sql/stmt"
)

// Logger :
type Logger interface {
	Debug(stmt *sqlstmt.Statement)
}
