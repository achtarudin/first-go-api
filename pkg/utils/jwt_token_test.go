package utils

import (
	"cutbray/first_api/domain/auth/entity"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type JwtTokenTestSuite struct {
	suite.Suite

	User        *entity.User
	TokenString string
	JwtToken    *jwt.Token
}

func (suite *JwtTokenTestSuite) SetupSuite() {
	suite.User = &entity.User{
		ID:    10,
		Name:  "Test User",
		Email: "userss@example.com",
	}
	// Setup code before each test
}

func (suite *JwtTokenTestSuite) TearDownTest() {
	// Cleanup code after each test
}

func (suite *JwtTokenTestSuite) TestGenerateJwtToken() {
	tokenString, err := GenerateToken(suite.User)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), tokenString)
}

func (suite *JwtTokenTestSuite) TestGenerateJwtTokenFromIdAndEmail() {
	tokenString, err := GenerateTokenFromIdAndEmail(suite.User.ID, suite.User.Email)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), tokenString)
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
