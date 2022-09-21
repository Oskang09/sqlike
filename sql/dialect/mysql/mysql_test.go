package mysql

import (
	"testing"

	sqlstmt "github.com/RevenueMonster/sqlike/sql/stmt"
	"github.com/stretchr/testify/require"
)

func TestGetVersion(t *testing.T) {
	ms := New()
	stmt := sqlstmt.AcquireStmt(ms)
	defer sqlstmt.ReleaseStmt(stmt)
	ms.GetVersion(stmt)
	require.Equal(t, "SELECT VERSION();", stmt.String())
	require.ElementsMatch(t, []interface{}{}, stmt.Args())
}
