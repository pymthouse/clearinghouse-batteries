package store

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestManifestCostUsesManifestIndex(t *testing.T) {
	ctx := context.Background()
	db, err := Open(ctx, filepath.Join(t.TempDir(), "plan.db"), true)
	require.NoError(t, err)
	defer db.Close()
	rows, err := db.DB.QueryContext(ctx, "EXPLAIN QUERY PLAN "+manifestCostQuery, "m-1")
	require.NoError(t, err)
	defer rows.Close()
	var plan []string
	for rows.Next() {
		var id, parent, unused int
		var detail string
		require.NoError(t, rows.Scan(&id, &parent, &unused, &detail))
		plan = append(plan, detail)
	}
	require.NoError(t, rows.Err())
	require.Contains(t, strings.Join(plan, "\n"), "SEARCH u USING INDEX usage_manifest (manifest_id=?)")
}
