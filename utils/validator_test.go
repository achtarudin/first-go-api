package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewValidator_Instance(t *testing.T) {
	assert.NotNil(t, NewValidator())
}

func TestValidateStruct_NotValid(t *testing.T) {

	// Arrange
	validator := NewValidator()

	// Act
	message, isValid := validator.ValidateStruct(struct {
		Email    string `validate:"required,email"`
		Password string `validate:"required,min=8"`
	}{
		Email:    "",
		Password: "invalid",
	})

	// Assert
	assert.NotNil(t, message)
	assert.Equal(t, isValid, false)

}

func TestValidateStruct_Valid(t *testing.T) {

	// Arrange
	validator := NewValidator()

	// Act
	message, isValid := validator.ValidateStruct(struct {
		Email    string `validate:"required,email"`
		Password string `validate:"required,min=8"`
	}{
		Email:    "example@email.com",
		Password: "password123",
	})

	// Assert
	assert.Nil(t, message)
	assert.Equal(t, isValid, true)

}

func TestValidateVar_NotValid(t *testing.T) {

	// Arrange
	validator := NewValidator()

	// Act
	isValid := validator.ValidateVar("invalid-email", "required,email")

	// Assert
	assert.NotNil(t, isValid)
	assert.Equal(t, isValid != nil, true)
}

func TestValidateVar_Valid(t *testing.T) {

	// Arrange
	validator := NewValidator()

	// Act
	isValid := validator.ValidateVar("example@email.com", "required,email")

	// Assert
	assert.Nil(t, isValid)
	assert.Equal(t, isValid == nil, true)
}
