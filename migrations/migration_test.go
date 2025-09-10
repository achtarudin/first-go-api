package migrations_test

import (
	"cutbray/first_api/infra"
	"cutbray/first_api/migrations"
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MigrationTestSuite struct {
	suite.Suite
	db *infra.Database
}

func (suite *MigrationTestSuite) SetupSuite() {

	viper.AutomaticEnv()

	viper.Set("DB_DATABASE", "first_go_api_test")
	viper.Set("DB_USER", "root")

	db, err := infra.NewDatabaseEnv()
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), db)

	suite.db = db
	assert.NotNil(suite.T(), suite.db)
}

func (suite *MigrationTestSuite) TestRunUpMigrations() {
	db, err := suite.db.DB.DB()

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), db)

	// run up migrations
	versi1, err := migrations.RunUpMigrations(db, viper.GetString("DB_DRIVER"))

	var totalTables1 int64
	suite.db.DB.Raw("SELECT count(*) FROM information_schema.tables WHERE table_schema = ?", viper.GetString("DB_DATABASE")).Scan(&totalTables1)
	fmt.Printf("Total totalTables1 di database: %d, version: %v\n", totalTables1, *versi1)

	assert.NoError(suite.T(), err)
	assert.True(suite.T(), *versi1 != 0)

	// run refresh the migrations
	versi2, err := migrations.RunRefreshMigrations(db, viper.GetString("DB_DRIVER"))

	var totalTables2 int64
	suite.db.DB.Raw("SELECT count(*) FROM information_schema.tables WHERE table_schema = ?", viper.GetString("DB_DATABASE")).Scan(&totalTables2)
	fmt.Printf("Total totalTables2 di database: %d, version: %v\n", totalTables2, *versi2)

	assert.NoError(suite.T(), err)
	assert.True(suite.T(), *versi2 != 0)
}

func TestMigration(t *testing.T) {
	suite.Run(t, new(MigrationTestSuite))
}
