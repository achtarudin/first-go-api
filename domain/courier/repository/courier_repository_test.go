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
	appConfig, db, err := migration.UsingFreshDatabaseTesting()
	assert.NoError(suite.T(), err)

	suite.db = db
	suite.config = appConfig

	assert.NotNil(suite.T(), suite.db)
	assert.NotNil(suite.T(), suite.config)

	// Initialize the repository
	suite.courierRepository = NewCourierRepository(suite.db)
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

func (suite *CourierRepositoryTestSuite) TestFindRoleCourierNotFound() {
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
	// return
	context := context.Background()
	err := suite.courierRepository.Trx(context, func(tx *gorm.DB) error {

		timestamp := time.Now().Unix()

		hashedPassword, err := utils.HashPassword("password")
		roleId, err := suite.courierRepository.FindRoleCourier(context, model.RoleCourier, tx)
		assert.NoError(suite.T(), err)

		courier := &entity.Courier{
			Name:     "John Doe",
			Email:    fmt.Sprintf("%d@email.com", timestamp),
			RoleId:   int(roleId),
			Password: hashedPassword,
			Phone:    fmt.Sprintf("+62%d", timestamp),
		}

		newCourier, err := suite.courierRepository.Create(context, courier, tx)
		assert.NoError(suite.T(), err)
		assert.NotZero(suite.T(), newCourier.ID)
		return err
	})
	assert.NoError(suite.T(), err)
}

func (suite *CourierRepositoryTestSuite) TestSaveCourierDuplicateEmail() {
	context := context.Background()
	err := suite.courierRepository.Trx(context, func(tx *gorm.DB) error {

		timestamp := time.Now().Unix()

		hashedPassword, err := utils.HashPassword("password")
		roleId, err := suite.courierRepository.FindRoleCourier(context, model.RoleCourier, tx)
		assert.NoError(suite.T(), err)

		courier1 := &entity.Courier{
			Name:      "John Doe",
			Email:     fmt.Sprintf("%d@email.com", timestamp),
			RoleId:    int(roleId),
			Password:  hashedPassword,
			Phone:     fmt.Sprintf("+62%d", timestamp),
			Latitude:  -6.200000,
			Longitude: 106.816666,
		}

		newCourier1, err := suite.courierRepository.Create(context, courier1, tx)
		assert.NoError(suite.T(), err)
		assert.NotZero(suite.T(), newCourier1.ID)

		courier2 := &entity.Courier{
			Name:     "John Doe",
			Email:    fmt.Sprintf("%d@email.com", timestamp),
			RoleId:   int(roleId),
			Password: hashedPassword,
			Phone:    fmt.Sprintf("+62%d", timestamp),
		}
		_, err = suite.courierRepository.Create(context, courier2, tx)

		return err
	})
	assert.Error(suite.T(), err)
}

func (suite *CourierRepositoryTestSuite) TestSaveCourierWithoutTx() {
	context := context.Background()
	err := suite.courierRepository.Trx(context, func(tx *gorm.DB) error {

		timestamp := time.Now().Unix()

		hashedPassword, err := utils.HashPassword("password")
		roleId, err := suite.courierRepository.FindRoleCourier(context, model.RoleCourier, nil)
		assert.NoError(suite.T(), err)

		courier := &entity.Courier{
			Name:     "John Doe",
			Email:    fmt.Sprintf("%d@email.com", timestamp),
			RoleId:   int(roleId),
			Password: hashedPassword,
			Phone:    fmt.Sprintf("+62%d", timestamp),
		}

		newCourier, err := suite.courierRepository.Create(context, courier, nil)
		assert.NoError(suite.T(), err)
		assert.NotZero(suite.T(), newCourier.ID)

		return err
	})
	assert.NoError(suite.T(), err)
}

func TestCourierRepository(t *testing.T) {
	suite.Run(t, new(CourierRepositoryTestSuite))
}
