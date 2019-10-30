package util

import (
	"chat/config"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type MyCustomClaims struct {
	jwt.StandardClaims
}

func CreateToken(id string) (string, error) {
	claims := MyCustomClaims{
		jwt.StandardClaims{
			Id:id,
			ExpiresAt: time.Now().AddDate(0,0,1).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.ENV_TOKEN_SECRET))
}

func ParseToken(tokenString string) (*MyCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.ENV_TOKEN_SECRET), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*MyCustomClaims)
	if ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("not a valid token")
}