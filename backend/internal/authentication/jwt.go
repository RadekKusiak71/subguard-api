package authentication

import (
	"strconv"
	"time"

	"github.com/RadekKusiak71/subguard-api/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

var (
	secretKey = []byte(config.Config.JWTSecret)
	tokenExp  = time.Second * time.Duration(config.Config.JWTExpireTime)
)

func CreateToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userID": strconv.Itoa(userID),
			"exp":    time.Now().Add(tokenExp).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	return token, nil
}
