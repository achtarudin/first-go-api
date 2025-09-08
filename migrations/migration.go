package migrations

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

func GetMigrationsFS(db *sql.DB) error {

	goose.SetBaseFS(embedMigrations)

	err := goose.SetDialect("sqlite3")

	if err != nil {
		return err
	}
	return goose.Up(db, ".")

}
