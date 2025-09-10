package migrations

import (
	"context"
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

/** Documents:
 * RunUpMigrations runs all up migrations
 * @Param db *sql.DB - the database connection
 * @Param dialect string - the database dialect (e.g., "mysql", "sqlite3")
 */
func RunUpMigrations(db *sql.DB, dialect string) (version *int64, err error) {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(dialect); err != nil {
		return nil, err
	}

	context := context.Background()

	if err := goose.UpContext(context, db, "."); err != nil {
		return nil, err
	}

	migrateVersion, err := goose.GetDBVersionContext(context, db)
	if err != nil {
		return nil, err
	}

	return &migrateVersion, nil
}

/** Documents:
 * RunRefreshMigrations rolls back all migrations and re-applies them.
 * @Param db *sql.DB - the database connection
 * @Param dialect string - the database dialect (e.g., "mysql", "sqlite3")
 */
func RunRefreshMigrations(db *sql.DB, dialect string) (version *int64, err error) {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(dialect); err != nil {
		return nil, err
	}

	context := context.Background()

	if err := goose.DownToContext(context, db, ".", 0); err != nil {
		return nil, err
	}

	if err := goose.UpContext(context, db, "."); err != nil {
		return nil, err
	}

	migrateVersion, err := goose.GetDBVersionContext(context, db)
	if err != nil {
		return nil, err
	}

	return &migrateVersion, nil

}
