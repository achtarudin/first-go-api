package utils

import (
	"cutbray/first_api/domain/auth/entity"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type JwtTokenTestSuite struct {
	suite.Suite
	User *entity.User
}

func (suite *JwtTokenTestSuite) SetupSuite() {
	suite.User = &entity.User{
		ID:    10,
		Name:  "Test User",
		Email: "users@example.com",
	}
}

func (suite *JwtTokenTestSuite) TearDownTest() {
	// Cleanup code after each test
}

func (suite *JwtTokenTestSuite) TestGenerateJwtTokenWithConfig() {

	jwtConfig := NewJwtConfig(suite.User, time.Now().Add(time.Minute).Unix())
	tokenString, err := GenerateTokenWithConfig(jwtConfig)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), tokenString)

	jwtToken, err := VerifyToken(tokenString)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), jwtToken)

	isValid, userMap, err := CheckTokenValid(jwtToken)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), isValid)
	assert.NotNil(suite.T(), userMap)
	assert.Equal(suite.T(), float64(suite.User.ID), userMap["user_id"])
	assert.Equal(suite.T(), suite.User.Email, userMap["email"])
}

func (suite *JwtTokenTestSuite) TestGenerateJwtToken() {
	tokenString, err := GenerateToken(suite.User)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), tokenString)

	jwtToken, err := VerifyToken(tokenString)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), jwtToken)

	isValid, userMap, err := CheckTokenValid(jwtToken)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), isValid)
	assert.NotNil(suite.T(), userMap)
	assert.Equal(suite.T(), float64(suite.User.ID), userMap["user_id"])
	assert.Equal(suite.T(), suite.User.Email, userMap["email"])
}

func (suite *JwtTokenTestSuite) TestGenerateJwtTokenFromIdAndEmail() {
	tokenString, err := GenerateTokenFromIdAndEmail(suite.User.ID, suite.User.Email)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), tokenString)

	jwtToken, err := VerifyToken(tokenString)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), jwtToken)

	isValid, userMap, err := CheckTokenValid(jwtToken)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), isValid)
	assert.NotNil(suite.T(), userMap)
	assert.Equal(suite.T(), float64(suite.User.ID), userMap["user_id"])
	assert.Equal(suite.T(), suite.User.Email, userMap["email"])
}

func (suite *JwtTokenTestSuite) TestVerifyJwtTokenExpiration() {
	jwtConfig := NewJwtConfig(suite.User, time.Now().Add(-time.Minute).Unix())
	tokenString, _ := GenerateTokenWithConfig(jwtConfig)
	jwtToken, err := VerifyToken(tokenString)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), jwtToken)
	assert.Equal(suite.T(), err, jwt.ErrTokenExpired)
}

func (suite *JwtTokenTestSuite) TestVerifyJwtToken() {
	tokenString, _ := GenerateToken(suite.User)
	jwtToken, err := VerifyToken(tokenString)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), jwtToken)
}

func (suite *JwtTokenTestSuite) TestCheckJwtTokenValid() {
	tokenString, _ := GenerateToken(suite.User)
	assert.NotEmpty(suite.T(), tokenString)

	jwtToken, _ := VerifyToken(tokenString)
	assert.NotNil(suite.T(), jwtToken)

	isValid, userMap, err := CheckTokenValid(jwtToken)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), isValid)
	assert.NotNil(suite.T(), userMap)
	assert.Equal(suite.T(), float64(suite.User.ID), userMap["user_id"])
	assert.Equal(suite.T(), suite.User.Email, userMap["email"])

}

func TestJwtToken(t *testing.T) {
	suite.Run(t, new(JwtTokenTestSuite))
}
