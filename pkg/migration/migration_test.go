package migration_test

import (
	"cutbray/first_api/infra"
	migration "cutbray/first_api/pkg/migration"
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MigrationTestSuite struct {
	suite.Suite
	db     *infra.Database
	config *viper.Viper
}

func (suite *MigrationTestSuite) SetupSuite() {

	appConfig, db, err := migration.UsingFreshDatabaseTesting()
	assert.NoError(suite.T(), err)

	suite.db = db
	suite.config = appConfig

	assert.NotNil(suite.T(), suite.db)
	assert.NotNil(suite.T(), suite.config)

}

func (suite *MigrationTestSuite) TearDownSuite() {
	if suite.db != nil {
		err := suite.db.Close()
		assert.NoError(suite.T(), err)
		fmt.Println("Database connection closed.")
	}

}

func (suite *MigrationTestSuite) TestRunUpMigrations() {
	db, err := suite.db.DB.DB()

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), db)

	// run up migration
	versi1, err := migration.RunUpMigrations(db, suite.config.GetString("DB_DRIVER"))
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), *versi1 != 0)

	var totalTables1 int64
	suite.db.DB.Raw("SELECT count(*) FROM information_schema.tables WHERE table_schema = ?", suite.config.GetString("DB_DATABASE_TESTING")).Scan(&totalTables1)
	fmt.Printf("Total totalTables1 di database: %d, version: %v\n", totalTables1, *versi1)

	// run refresh the migration
	versi2, err := migration.RunRefreshMigrations(db, suite.config.GetString("DB_DRIVER"))
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), *versi2 != 0)

	var totalTables2 int64
	suite.db.DB.Raw("SELECT count(*) FROM information_schema.tables WHERE table_schema = ?", suite.config.GetString("DB_DATABASE_TESTING")).Scan(&totalTables2)
	fmt.Printf("Total totalTables2 di database: %d, version: %v\n", totalTables2, *versi2)

}

func TestMigration(t *testing.T) {
	suite.Run(t, new(MigrationTestSuite))
}
