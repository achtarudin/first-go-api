package repository

import (
	"context"
	"cutbray/first_api/domain/courier/entity"
	"cutbray/first_api/infra"
	"cutbray/first_api/pkg/migration"
	"cutbray/first_api/pkg/model"
	"cutbray/first_api/pkg/utils"
	"fmt"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type CourierRepositoryTestSuite struct {
	suite.Suite
	db                *infra.Database
	config            *viper.Viper
	courierRepository CourierRepository
}

func (suite *CourierRepositoryTestSuite) SetupSuite() {
	config := infra.NewAppConfig()
	assert.NotNil(suite.T(), config)

	err := config.LoadEnvConfig(nil)
	assert.NoError(suite.T(), err)

	suite.config = config.GetViper()
	assert.NotNil(suite.T(), suite.config)

	// Config database
	db, err := infra.NewDatabase(infra.DatabaseConfig{
		Host:     suite.config.GetString("DB_HOST"),
		Port:     suite.config.GetInt("DB_PORT"),
		User:     suite.config.GetString("DB_USER"),
		Password: suite.config.GetString("DB_PASSWORD"),
		DBName:   suite.config.GetString("DB_DATABASE_TESTING"),
	})
	suite.db = db
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), suite.db)

	sqlDB, err := suite.db.DB.DB()
	_, err = migration.RunUpMigrations(sqlDB, suite.config.GetString("DB_DRIVER"))
	assert.NoError(suite.T(), err)

	suite.courierRepository = NewCourierRepository(suite.db.DB)
	assert.NotNil(suite.T(), suite.courierRepository)
}

func (suite *CourierRepositoryTestSuite) TearDownSuite() {
	if suite.db != nil {
		err := suite.db.Close()
		assert.NoError(suite.T(), err)
		fmt.Println("Database connection closed.")
	}
}

func (suite *CourierRepositoryTestSuite) TestFindRoleCourier() {
	context := context.Background()
	err := suite.courierRepository.Trx(context, func(tx *gorm.DB) error {
		roleId, err := suite.courierRepository.FindRoleCourier(context, model.RoleCourier, tx)
		assert.NoError(suite.T(), err)
		assert.NotZero(suite.T(), roleId)
		return err
	})

	assert.NoError(suite.T(), err)

}

func (suite *CourierRepositoryTestSuite) TestNotFoundFindRoleCourier() {
	context := context.Background()
	err := suite.courierRepository.Trx(context, func(tx *gorm.DB) error {
		roleId, err := suite.courierRepository.FindRoleCourier(context, model.RoleNotFound, tx)
		assert.Error(suite.T(), err)
		assert.Zero(suite.T(), roleId)
		return err
	})
	assert.Error(suite.T(), err)
}

func (suite *CourierRepositoryTestSuite) TestSaveCourier() {
	context := context.Background()
	err := suite.courierRepository.Trx(context, func(tx *gorm.DB) error {

		timestamp := time.Now().Unix()

		hashedPassword, err := utils.HashPassword("password")
		roleId, err := suite.courierRepository.FindRoleCourier(context, model.RoleCourier, tx)
		assert.NoError(suite.T(), err)

		newCourier := &entity.Courier{
			Name:     "John Doe",
			Email:    fmt.Sprintf("%d@email.com", timestamp),
			RoleID:   int(roleId),
			Password: hashedPassword,
			Phone:    fmt.Sprintf("+62%d", timestamp),
		}

		err = suite.courierRepository.Create(context, newCourier, tx)
		assert.NoError(suite.T(), err)
		assert.NotZero(suite.T(), newCourier.ID)
		return err
	})
	assert.NoError(suite.T(), err)
}

func (suite *CourierRepositoryTestSuite) TestErroSaveCourier() {
	context := context.Background()
	err := suite.courierRepository.Trx(context, func(tx *gorm.DB) error {

		timestamp := time.Now().Unix()

		hashedPassword, err := utils.HashPassword("password")
		roleId, err := suite.courierRepository.FindRoleCourier(context, model.RoleCourier, tx)
		assert.NoError(suite.T(), err)

		newCourier := &entity.Courier{
			Name:     "John Doe",
			Email:    fmt.Sprintf("%d@email.com", timestamp),
			RoleID:   int(roleId),
			Password: hashedPassword,
			Phone:    fmt.Sprintf("+62%d", timestamp),
		}

		err = suite.courierRepository.Create(context, newCourier, tx)

		assert.NoError(suite.T(), err)
		assert.NotZero(suite.T(), newCourier.ID)
		return err
	})
	assert.NoError(suite.T(), err)
}

func TestCourierRepository(t *testing.T) {
	suite.Run(t, new(CourierRepositoryTestSuite))
}
