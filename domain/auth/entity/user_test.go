package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserEntity(t *testing.T) {

	user := NewUser(1, "Test User", "test@example.com", "password")

	assert.NotNil(t, user)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "password", user.Password)
	assert.Empty(t, user.Token)
}
