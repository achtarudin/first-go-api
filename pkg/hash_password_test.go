package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "securePassword123!"
	wrongPassword := "wrongPassword!"

	var passwordHashed string
	t.Run("should hash password successfully", func(t *testing.T) {

		// Act
		hash, err := HashPassword(password)
		passwordHashed = hash

		// Asserts
		assert.NoError(t, err)
		assert.NotEqual(t, password, hash)
	})

	t.Run("should verify correct password", func(t *testing.T) {
		// Act
		isValid := VerifyPassword(password, passwordHashed)

		// Asserts
		assert.NotNil(t, passwordHashed)
		assert.True(t, isValid)
	})

	t.Run("should fail to verify incorrect password", func(t *testing.T) {

		// Act
		isValid := VerifyPassword(wrongPassword, passwordHashed)

		// Asserts
		assert.NotNil(t, passwordHashed)
		assert.False(t, isValid)

	})
}
