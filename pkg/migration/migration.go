package migration

import (
	"context"
	"cutbray/first_api/infra"
	"database/sql"
	"embed"
	"errors"

	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
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

func UsingFreshDatabaseTesting() (appConfig *viper.Viper, database *infra.Database, err error) {
	config := infra.NewAppConfig()

	err = config.LoadEnvConfig(nil)

	if err != nil {
		return nil, nil, errors.New("failed to load env config: " + err.Error())
	}

	viperConfig := config.GetViper()

	// Config database
	db, err := infra.NewDatabase(infra.DatabaseConfig{
		Host:     viperConfig.GetString("DB_HOST"),
		Port:     viperConfig.GetInt("DB_PORT"),
		User:     viperConfig.GetString("DB_USER"),
		Password: viperConfig.GetString("DB_PASSWORD"),
		DBName:   viperConfig.GetString("DB_DATABASE_TESTING"),
	})

	sqlDB, err := db.DB.DB()

	if err != nil {
		return nil, nil, errors.New("failed to get sql.DB from gorm.DB: " + err.Error())
	}

	_, err = RunRefreshMigrations(sqlDB, viperConfig.GetString("DB_DRIVER"))

	if err != nil {
		return nil, nil, errors.New("failed to run refresh migrations: " + err.Error())
	}

	return viperConfig, db, nil

}
