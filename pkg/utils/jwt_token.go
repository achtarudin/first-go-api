package utils

import (
	"cutbray/first_api/domain/auth/entity"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var HmacSampleSecret = []byte("8MN1DHO6SI.EOA75KHT2C.FQO08FA9IH")

type jwtConfig struct {
	user *entity.User
	exp  int64
}

func NewJwtConfig(user *entity.User, exp int64) *jwtConfig {
	return &jwtConfig{
		user: user,
		exp:  exp,
	}
}

func GenerateTokenWithConfig(config *jwtConfig) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": config.user.ID,
		"email":   config.user.Email,
		"exp":     config.exp,
		"iat":     time.Now().Unix(),
	})
	return token.SignedString(HmacSampleSecret)
}

func GenerateToken(user *entity.User) (string, error) {
	jwtConfig := NewJwtConfig(user, time.Now().Add(time.Hour*72).Unix())
	return GenerateTokenWithConfig(jwtConfig)
}

func GenerateTokenFromIdAndEmail(id int, email string) (string, error) {

	userEntity := &entity.User{
		ID:    id,
		Email: email,
	}

	jwtConfig := NewJwtConfig(userEntity, time.Now().Add(time.Hour*72).Unix())

	return GenerateTokenWithConfig(jwtConfig)

}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return HmacSampleSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, jwt.ErrTokenExpired
		}
		return nil, err
	}

	return token, nil
}

func CheckTokenValid(token *jwt.Token) (isValid bool, userMap map[string]interface{}, err error) {
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return false, nil, fmt.Errorf("invalid token claims")
	}

	// Ambil user_id dan email dari claims
	userId, idExists := claims["user_id"]
	email, emailExists := claims["email"]

	if !idExists || !emailExists {
		return false, nil, fmt.Errorf("invalid token claims")
	}

	userData := map[string]interface{}{}

	userData["user_id"] = userId
	userData["email"] = email

	return true, userData, nil
}
