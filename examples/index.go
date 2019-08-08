package examples

import (
	"testing"

	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/indexes"
	"github.com/stretchr/testify/require"
)

// IndexExamples :
func IndexExamples(t *testing.T, db *sqlike.Database) {
	var (
		err  error
		idxs []sqlike.Index
	)

	table := db.Table("Index")

	{

		err = table.DropIfExits()
		require.NoError(t, err)
	}

	// Migrate and create unique index with `unique_index` tag
	{
		err = table.Migrate(indexStruct{})
		require.NoError(t, err)
	}

	// Create one index
	{
		idx := table.Indexes()
		err = idx.CreateOne(indexes.Index{
			Columns: []indexes.Column{
				indexes.Column{Name: "ID"},
			},
		})
		require.NoError(t, err)
		idxs, err = idx.List()
		require.True(t, len(idxs) > 1)
	}

	{
		err = db.BuildIndexes()
		require.NoError(t, err)
		idxs, err = db.Table("A").Indexes().List()
		require.NoError(t, err)
		require.Contains(t, idxs, sqlike.Index{
			Name:      "IX_SID_ASC-IX_Emoji_ASC-IX_Bool_DESC",
			Type:      "BTREE",
			IsVisible: true,
		})
		require.Contains(t, idxs, sqlike.Index{
			Name:      "test_idx",
			Type:      "BTREE",
			IsVisible: true,
		})
	}

}