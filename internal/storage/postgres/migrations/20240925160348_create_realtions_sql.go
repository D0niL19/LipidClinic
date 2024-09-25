package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateRealtionsSql, downCreateRealtionsSql)
}

func upCreateRealtionsSql(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return nil
}

func downCreateRealtionsSql(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
