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

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return HmacSampleSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return nil, err
	}

	// if claims, ok := token.Claims.(jwt.MapClaims); ok {
	// 	fmt.Println(claims["user_id"], claims["email"])
	// } else {
	// 	fmt.Println(err)
	// }

	return token, nil
}

func CheckTokenValid(token *jwt.Token) (bool, error) {
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return false, fmt.Errorf("invalid token claims")
	}

	// Ambil exp dan iat
	exp := int64(claims["exp"].(float64))
	iat := int64(claims["iat"].(float64))

	now := time.Now().Unix()

	return exp > now && iat < now, nil
}
