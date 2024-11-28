package dialect

import (
	"reflect"
	"strings"
	"sync"

	"github.com/Oskang09/sqlike/reflext"
	"github.com/Oskang09/sqlike/sql"
	"github.com/Oskang09/sqlike/sql/codec"
	"github.com/Oskang09/sqlike/sql/driver"
	sqlstmt "github.com/Oskang09/sqlike/sql/stmt"
	"github.com/Oskang09/sqlike/sql/util"
	"github.com/Oskang09/sqlike/sqlike/actions"
	"github.com/Oskang09/sqlike/sqlike/indexes"
	"github.com/Oskang09/sqlike/sqlike/options"
)

// SQLDialect :
type SQLDialect interface {
	TableName(db, table string) string
	Var(i int) string
	Quote(n string) string
	Format(v interface{}) (val string)
}

// Dialect :
type Dialect interface {
	SQLDialect
	Connect(opt *options.ConnectOptions) (connStr string)
	UseDatabase(stmt sqlstmt.Stmt, db string)
	GetVersion(stmt sqlstmt.Stmt)
	GetDatabases(stmt sqlstmt.Stmt)
	CreateDatabase(stmt sqlstmt.Stmt, db string, checkExists bool)
	DropDatabase(stmt sqlstmt.Stmt, db string, checkExists bool)
	HasTable(stmt sqlstmt.Stmt, db, table string)
	HasPrimaryKey(stmt sqlstmt.Stmt, db, table string)
	RenameTable(stmt sqlstmt.Stmt, db, oldName, newName string)
	RenameColumn(stmt sqlstmt.Stmt, db, table, oldColName, newColName string)
	DropColumn(stmt sqlstmt.Stmt, db, table, column string)
	DropTable(stmt sqlstmt.Stmt, db, table string, checkExists bool)
	TruncateTable(stmt sqlstmt.Stmt, db, table string)
	GetColumns(stmt sqlstmt.Stmt, db, table string)
	HasIndexByName(stmt sqlstmt.Stmt, db, table, indexName string)
	HasIndex(stmt sqlstmt.Stmt, dbName, table string, idx indexes.Index)
	GetIndexes(stmt sqlstmt.Stmt, db, table string)
	CreateIndexes(stmt sqlstmt.Stmt, db, table string, idxs []indexes.Index, supportDesc bool)
	DropIndexes(stmt sqlstmt.Stmt, db, table string, idxs []string)
	CreateTable(stmt sqlstmt.Stmt, db, table, pk string, info driver.Info, fields []reflext.StructFielder) (err error)
	AlterTable(stmt sqlstmt.Stmt, db, table, pk string, hasPk bool, info driver.Info, fields []reflext.StructFielder, columns util.StringSlice, indexes util.StringSlice, unsafe bool) (err error)
	InsertInto(stmt sqlstmt.Stmt, db, table, pk string, mapper reflext.StructMapper, codec codec.Codecer, fields []reflext.StructFielder, values reflect.Value, opts *options.InsertOptions) (err error)
	Select(stmt sqlstmt.Stmt, act *actions.FindActions, mode options.LockMode) (err error)
	Update(stmt sqlstmt.Stmt, act *actions.UpdateActions) (err error)
	Delete(stmt sqlstmt.Stmt, act *actions.DeleteActions) (err error)
	SelectStmt(stmt sqlstmt.Stmt, query interface{}) (err error)
	Replace(stmt sqlstmt.Stmt, db, table string, columns []string, query *sql.SelectStmt) (err error)
}

var (
	mutex    = new(sync.RWMutex)
	dialects = make(map[string]Dialect)
)

// RegisterDialect :
func RegisterDialect(driver string, dialect Dialect) {
	mutex.Lock()
	defer mutex.Unlock()
	if dialect == nil {
		panic("invalid nil dialect")
	}
	dialects[driver] = dialect
}

// GetDialectByDriver :
func GetDialectByDriver(driver string) Dialect {
	mutex.RLock()
	defer mutex.RUnlock()
	driver = strings.TrimSpace(strings.ToLower(driver))
	return dialects[driver]
}
