package token

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

const (
	userPhone = "uphone"
)

type TokenManager struct {
	signingkey string
}

func NewManager(skey string) *TokenManager {
	return &TokenManager{
		signingkey: skey,
	}
}

func (t *TokenManager) Create(phone string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims[userPhone] = phone

	return token.SignedString([]byte(t.signingkey))
}

func (t *TokenManager) Parse(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(t.signingkey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error get user claims from token")
	}

	phone := claims[userPhone].(string)

	return phone, nil
}
