package utils

import (
	"cutbray/first_api/domain/auth/entity"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var HmacSampleSecret = []byte("your-very-secret-key")

func GenerateToken(user *entity.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
		"iat":     time.Now().Unix(),
		// "nbf":     time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})
	return token.SignedString(HmacSampleSecret)
}

func GenerateTokenFromIdAndEmail(id int, email string) (string, error) {
	return GenerateToken(&entity.User{
		ID:    id,
		Email: email,
	})
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return HmacSampleSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return nil, err
	}

	return token, nil
}

func CheckTokenValid(token *jwt.Token) (isValid bool, userMap map[string]interface{}, err error) {
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return false, nil, fmt.Errorf("invalid token claims")
	}

	expResult, err := claims.GetExpirationTime()
	if err != nil {
		return false, nil, fmt.Errorf("invalid exp token claims")
	}

	iatResult, err := claims.GetIssuedAt()
	if err != nil {
		return false, nil, fmt.Errorf("invalid iat token claims")
	}

	// Ambil user_id dan email dari claims
	userData := map[string]interface{}{}

	if val, ok := claims["user_id"]; ok {
		userData["user_id"] = val
	}

	if val, ok := claims["email"]; ok {
		userData["email"] = val
	}

	exp := expResult.Unix()
	iat := iatResult.Unix()
	now := time.Now().Unix()

	return exp > now && iat <= now, userData, nil
}
