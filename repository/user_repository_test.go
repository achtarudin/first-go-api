package repository

import (
	"cutbray/first_api/infra"
	"cutbray/first_api/repository/model"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestCreate(t *testing.T) {
	mockRepo := new(MockUserRepository)
	user := &model.User{Username: "testuser", Email: "testuser@example.com"}
	mockRepo.On("Create", user).Return(nil)

	err := mockRepo.Create(user)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreate2(t *testing.T) {
	viper.AutomaticEnv()
	db, err := infra.NewDatabaseEnv()

	assert.NoError(t, err)
	assert.NotNil(t, db)
	repo := NewUserRepository(db.DB)

	user := &model.User{Username: "testuser3sd", Email: "testuser2@exampls3e.com", Password: "password"}
	err = repo.Create(user)
	assert.NoError(t, err)
}
