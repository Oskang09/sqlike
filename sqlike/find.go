package sqlike

import (
	"context"
	"database/sql"

	"github.com/Oskang09/sqlike/reflext"
	"github.com/Oskang09/sqlike/sql/codec"
	sqldialect "github.com/Oskang09/sqlike/sql/dialect"
	sqldriver "github.com/Oskang09/sqlike/sql/driver"
	sqlstmt "github.com/Oskang09/sqlike/sql/stmt"
	"github.com/Oskang09/sqlike/sqlike/actions"
	"github.com/Oskang09/sqlike/sqlike/logs"
	"github.com/Oskang09/sqlike/sqlike/options"
	"github.com/Oskang09/sqlike/sqlike/primitive"
)

// SingleResult : single result is an interface implementing apis as similar as driver.Result
type SingleResult interface {
	Scan(dest ...interface{}) error
	Decode(dest interface{}) error
	Columns() []string
	ColumnTypes() ([]*sql.ColumnType, error)
	Error() error
}

// FindOne : find single record on the table, you should alway check the return error to ensure it have result return.
func (tb *Table) FindOne(ctx context.Context, act actions.SelectOneStatement, opts ...*options.FindOneOptions) SingleResult {
	x := new(actions.FindOneActions)
	if act != nil {
		*x = *(act.(*actions.FindOneActions))
	}
	opt := new(options.FindOneOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}
	x.Limit(1)
	rslt := find(
		ctx,
		tb.dbName,
		tb.name,
		tb.client.cache,
		tb.codec,
		tb.driver,
		tb.dialect,
		tb.logger,
		&x.FindActions,
		&opt.FindOptions,
		opt.FindOptions.LockMode,
	)
	rslt.close = true
	if rslt.err != nil {
		return rslt
	}
	if !rslt.Next() {
		rslt.err = sql.ErrNoRows
	}
	return rslt
}

// Find : find multiple records on the table.
func (tb *Table) Find(ctx context.Context, act actions.SelectStatement, opts ...*options.FindOptions) (*Result, error) {
	x := new(actions.FindActions)
	if act != nil {
		*x = *(act.(*actions.FindActions))
	}
	opt := new(options.FindOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}
	// has limit and limit value is zero
	if !opt.NoLimit && x.Count < 1 {
		x.Limit(100)
	}
	csr := find(
		ctx,
		tb.dbName,
		tb.name,
		tb.client.cache,
		tb.codec,
		tb.driver,
		tb.dialect,
		tb.logger,
		x,
		opt,
		opt.LockMode,
	)
	if csr.err != nil {
		return nil, csr.err
	}
	return csr, nil
}

func find(ctx context.Context, dbName, tbName string, cache reflext.StructMapper, cdc codec.Codecer, driver sqldriver.Driver, dialect sqldialect.Dialect, logger logs.Logger, act *actions.FindActions, opt *options.FindOptions, lock options.LockMode) *Result {
	if act.Database == "" {
		act.Database = dbName
	}
	if act.Table == "" {
		act.Table = tbName
	}

	groups := extractResolution(ctx)
	if len(groups) > 0 {
		if len(act.Conditions.Values) > 0 {
			act.Conditions.Values = append(act.Conditions.Values, primitive.And)
		}

		for _, group := range groups {
			act.Conditions.Values = append(act.Conditions.Values, group.Values...)
		}
	}

	rslt := new(Result)
	rslt.cache = cache
	rslt.codec = cdc

	stmt := sqlstmt.AcquireStmt(dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	if err := dialect.Select(stmt, act, lock); err != nil {
		rslt.err = err
		return rslt
	}
	rows, err := sqldriver.Query(
		ctx,
		driver,
		stmt,
		getLogger(logger, opt.Debug),
	)
	if err != nil {
		rslt.err = err
		return rslt
	}
	rslt.rows = rows
	rslt.columnTypes, rslt.err = rows.ColumnTypes()
	if rslt.err != nil {
		defer rslt.rows.Close()
	}
	for _, col := range rslt.columnTypes {
		rslt.columns = append(rslt.columns, col.Name())
	}
	return rslt
}
