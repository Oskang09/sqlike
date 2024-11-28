package mysql

import (
	"github.com/Oskang09/sqlike/sql/codec"
	"github.com/Oskang09/sqlike/sql/dialect"
	"github.com/Oskang09/sqlike/sql/schema"
	sqlstmt "github.com/Oskang09/sqlike/sql/stmt"
	sqlutil "github.com/Oskang09/sqlike/sql/util"
)

// MySQL :
type MySQL struct {
	schema *schema.Builder
	parser *sqlstmt.StatementBuilder
	sqlutil.MySQLUtil
}

var _ dialect.Dialect = (*(MySQL))(nil)

// New :
func New() *MySQL {
	sb := schema.NewBuilder()
	pr := sqlstmt.NewStatementBuilder()

	mySQLSchema{}.SetBuilders(sb)
	mySQLBuilder{}.SetRegistryAndBuilders(codec.DefaultRegistry, pr)

	return &MySQL{
		schema: sb,
		parser: pr,
	}
}

// GetVersion :
func (ms MySQL) GetVersion(stmt sqlstmt.Stmt) {
	stmt.WriteString("SELECT VERSION();")
}
