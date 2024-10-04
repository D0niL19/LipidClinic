package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(userId int64, userEmail string, duration time.Duration, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = userId
	claims["email"] = userEmail
	claims["exp"] = time.Now().Add(duration).Unix()

	accessToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return accessToken, nil

}
