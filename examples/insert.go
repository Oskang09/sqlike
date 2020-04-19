package examples

import (
	"context"
	"database/sql"
	"testing"

	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/stretchr/testify/require"
)

// InsertExamples :
func InsertExamples(t *testing.T, ctx context.Context, db *sqlike.Database) {
	var (
		err      error
		result   sql.Result
		affected int64
	)

	table := db.Table("NormalStruct")
	ns := newNormalStruct()

	// Single insert
	{
		result, err = table.InsertOne(
			ctx,
			&ns,
			options.InsertOne().
				SetOmitFields("Int").
				SetDebug(true))
		require.NoError(t, err)
		affected, err = result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(1), affected)
	}

	// Single upsert
	// - https://dev.mysql.com/doc/refman/8.0/en/insert-on-duplicate.html
	{
		ns.Emoji = `🤕`
		m := make(map[string]int)
		m["one"] = 1
		m["two"] = 2
		ns.Map = m
		result, err = table.InsertOne(
			ctx,
			&ns,
			options.InsertOne().
				SetDebug(true).
				SetMode(options.InsertOnDuplicate))
		require.NoError(t, err)
		affected, err = result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(2), affected)
	}

	// Multiple insert
	{
		nss := [...]normalStruct{
			newNormalStruct(),
			newNormalStruct(),
			newNormalStruct(),
		}
		result, err = table.Insert(
			ctx,
			&nss,
			options.Insert().SetDebug(true),
		)
		require.NoError(t, err)
		affected, err = result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(3), affected)
	}

	table2 := db.Table("PtrStruct")

	// Pointer insertion
	{
		_, err = table2.InsertOne(
			ctx,
			&ptrStruct{},
			options.InsertOne().SetDebug(true),
		)
		require.NoError(t, err)
	}

	// Pointer insertion
	{
		ps := []ptrStruct{
			newPtrStruct(),
			newPtrStruct(),
			newPtrStruct(),
			newPtrStruct(),
			newPtrStruct(),
		}
		_, err = table2.Insert(
			ctx,
			&ps,
			options.Insert().SetDebug(true),
		)
		require.NoError(t, err)
	}

	// Error insertion
	{
		_, err = table.InsertOne(
			ctx,
			&struct {
				Interface interface{}
			}{},
		)
		require.Error(t, err)
		_, err = table.InsertOne(ctx, struct{}{})
		require.Error(t, err)
		var empty *struct{}
		_, err = table.InsertOne(ctx, empty)
		require.Error(t, err)

		_, err = table.Insert(ctx, []interface{}{})
		require.Error(t, err)
	}

	table3 := db.Table("GeneratedStruct")

	// generated column insertion
	{
		cols := []*generatedStruct{
			newGeneratedStruct(),
			newGeneratedStruct(),
			newGeneratedStruct(),
			newGeneratedStruct(),
			newGeneratedStruct(),
		}
		_, err = table3.Insert(
			ctx,
			&cols,
			options.Insert().SetDebug(true),
		)
		require.NoError(t, err)
	}
}

// InsertErrorExamples :
func InsertErrorExamples(t *testing.T, db *sqlike.Database) {
	var (
		ns  normalStruct
		ctx = context.Background()
		err error
	)

	{
		_, err = db.Table("NormalStruct").InsertOne(ctx, nil)
		require.Error(t, err)

		var uninitialized *normalStruct
		_, err = db.Table("NormalStruct").InsertOne(ctx, uninitialized)
		require.Error(t, err)

		ns = normalStruct{}
		_, err = db.Table("NormalStruct").InsertOne(ctx, ns)
		require.Error(t, err)
	}

	{
		_, err = db.Table("NormalStruct").Insert(
			ctx,
			[]normalStruct{},
		)
		require.Error(t, err)
	}
}
