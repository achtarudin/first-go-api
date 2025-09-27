package repository

import (
	"context"
	"cutbray/first_api/domain/auth/entity"
	"cutbray/first_api/infra"
	"cutbray/first_api/pkg/customerror"
	"cutbray/first_api/pkg/migration"
	"cutbray/first_api/pkg/utils"
	"errors"
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AuthRepositoryTestSuite struct {
	suite.Suite
	db             *infra.Database
	config         *viper.Viper
	authRepository AuthRepository
}

func (suite *AuthRepositoryTestSuite) SetupSuite() {
	appConfig, db, err := migration.UsingFreshDatabaseTesting()
	assert.NoError(suite.T(), err)

	suite.db = db
	suite.config = appConfig

	assert.NotNil(suite.T(), suite.db)
	assert.NotNil(suite.T(), suite.config)

	// Initialize the repository
	suite.authRepository = NewAuthRepository(suite.db)
	assert.NotNil(suite.T(), suite.authRepository)
}

func (suite *AuthRepositoryTestSuite) TearDownSuite() {
	if suite.db != nil {
		err := suite.db.Close()
		assert.NoError(suite.T(), err)
		fmt.Println("Database connection closed.")
	}
}

func (suite *AuthRepositoryTestSuite) TestTrx_Success() {
	context := context.Background()
	err := suite.authRepository.Trx(context, func(tx *gorm.DB) error {
		return nil
	})

	assert.NoError(suite.T(), err)
}

func (suite *AuthRepositoryTestSuite) TestTrx_Failed() {
	var context = context.Background()
	var ErrTransactionFailed = errors.New("transaction failed")

	err := suite.authRepository.Trx(context, func(tx *gorm.DB) error {
		return fmt.Errorf("transaction failed: %w", ErrTransactionFailed)
	})

	assert.ErrorIs(suite.T(), err, ErrTransactionFailed)
	assert.Error(suite.T(), err)
}

func (suite *AuthRepositoryTestSuite) TestFindByEmail_Success() {

	var context = context.Background()

	foundUser := &entity.User{}

	newUser := &entity.User{
		Name:     "Test User 1",
		Email:    "user1@email.com",
		Password: "hashedpassword",
	}

	// Hash the password before saving
	hashedPassword, err := utils.HashPassword(newUser.Password)
	assert.NoError(suite.T(), err, "HashPassword should not return an error")

	err = suite.authRepository.Save(context, newUser, hashedPassword)
	assert.NoError(suite.T(), err, "Save should not return an error")

	foundUser, err = suite.authRepository.FindByEmail(context, newUser.Email)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), foundUser.Password)
	assert.NotEmpty(suite.T(), foundUser.ID)
	assert.NotEmpty(suite.T(), foundUser.Email)
}

func (suite *AuthRepositoryTestSuite) TestFindByEmail_Failed() {

	var context = context.Background()
	var appError = &customerror.CustomError{}

	_, err := suite.authRepository.FindByEmail(context, "<email>")

	assert.ErrorAs(suite.T(), err, &appError)
}

func (suite *AuthRepositoryTestSuite) TestSave_Success() {

	var context = context.Background()

	newUser := &entity.User{
		Name:     "Test User 2",
		Email:    "user2@email.com",
		Password: "hashedpassword",
	}

	// Hash the password before saving
	hashedPassword, err := utils.HashPassword(newUser.Password)
	assert.NoError(suite.T(), err, "HashPassword should not return an error")
	newUser.Password = hashedPassword

	err = suite.authRepository.Save(context, newUser, hashedPassword)
	assert.NoError(suite.T(), err)
}

func (suite *AuthRepositoryTestSuite) TestSave_Failed() {

	var context = context.Background()

	newUser := &entity.User{
		Name:     "Test User 3",
		Email:    "user3@email.com",
		Password: "hashedpassword",
	}

	// Hash the password before saving
	hashedPassword, err := utils.HashPassword(newUser.Password)
	assert.NoError(suite.T(), err, "HashPassword should not return an error")
	newUser.Password = hashedPassword

	// First save should succeed
	suite.authRepository.Save(context, newUser, hashedPassword)

	// Second save with the same email should fail
	err2 := suite.authRepository.Save(context, newUser, hashedPassword)

	appError := &customerror.CustomError{}
	assert.ErrorAs(suite.T(), err2, &appError)

}

func (suite *AuthRepositoryTestSuite) TestSaveTrx_Success() {

	var context = context.Background()

	newUser := &entity.User{
		Name:     "Test User 4",
		Email:    "user4@email.com",
		Password: "hashedpassword",
	}

	// Hash the password before saving
	hashedPassword, err := utils.HashPassword(newUser.Password)
	assert.NoError(suite.T(), err, "HashPassword should not return an error")
	newUser.Password = hashedPassword

	// Save the user within a transaction
	err = suite.authRepository.Trx(context, func(tx *gorm.DB) error {
		return suite.authRepository.SaveWithTrx(context, newUser, hashedPassword, tx)
	})
	assert.NoError(suite.T(), err)

	// Verify the user was saved
	foundUser, err := suite.authRepository.FindByEmail(context, newUser.Email)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), foundUser)

	// Additionally, verify by ID
	foundById, err := suite.authRepository.FindById(context, fmt.Sprintf("%d", newUser.ID))
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), foundById)

}

func (suite *AuthRepositoryTestSuite) TestReadAll_Success() {
	var context = context.Background()
	_, err := suite.authRepository.ReadAll(context)
	assert.NoError(suite.T(), err)
}

func (suite *AuthRepositoryTestSuite) TestSaveTrx_Failed() {

	context := context.Background()

	newUser := &entity.User{
		Name:     "Test User 5",
		Email:    "user5@email.com",
		Password: "hashedpassword",
	}

	// Hash the password before saving
	hashedPassword, err := utils.HashPassword(newUser.Password)
	assert.NoError(suite.T(), err, "HashPassword should not return an error")
	newUser.Password = hashedPassword

	// Save the user within a transaction and intentionally return an error to trigger a rollback
	errSave := suite.authRepository.Trx(context, func(tx *gorm.DB) error {
		suite.authRepository.SaveWithTrx(context, newUser, hashedPassword, tx)

		return customerror.New(customerror.CodeInternal, "Intentional Error to rollback",
			nil, errors.New("intentional error"),
		)
	})

	// Find the user to ensure it was not saved due to rollback
	foundUser, errFoundUser := suite.authRepository.FindByEmail(context, newUser.Email)

	appError := &customerror.CustomError{}
	assert.ErrorAs(suite.T(), errSave, &appError)
	assert.ErrorAs(suite.T(), errFoundUser, &appError)
	assert.Nil(suite.T(), foundUser)

}
func TestAuthRepositoryTest(t *testing.T) {
	suite.Run(t, new(AuthRepositoryTestSuite))
}
