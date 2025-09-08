package migrations

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

func RefreshSQLiteMigrations(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}
	if err := goose.DownTo(db, ".", 0, goose.WithNoVersioning()); err != nil {
		return err
	}
	if err := goose.Up(db, ".", goose.WithNoVersioning()); err != nil {
		return err
	}
	return nil
}

func RefreshMysqlMigrations(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("mysql"); err != nil {
		return err
	}
	if err := goose.DownTo(db, ".", 0, goose.WithNoVersioning()); err != nil {
		return err
	}
	if err := goose.Up(db, ".", goose.WithNoVersioning()); err != nil {
		return err
	}
	return nil
}
