package migrations

import (
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
func RunUpMigrations(db *sql.DB, dialect string) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(dialect); err != nil {
		return err
	}
	if err := goose.Up(db, "."); err != nil {
		return err
	}
	return nil
}

/** Documents:
 * RunRefreshMigrations rolls back all migrations and re-applies them.
 * @Param db *sql.DB - the database connection
 * @Param dialect string - the database dialect (e.g., "mysql", "sqlite3")
 */
func RunRefreshMigrations(db *sql.DB, dialect string) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(dialect); err != nil {
		return err
	}
	if err := goose.DownTo(db, ".", 0, goose.WithNoVersioning()); err != nil {
		return err
	}

	if err := goose.Up(db, "."); err != nil {
		return err
	}
	return nil
}
