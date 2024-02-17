package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var tokenLifetime = time.Duration(15) * time.Minute

type AuthManager struct {
	signingKey []byte
}

func NewAuthManager(signingKey []byte) *AuthManager {
	return &AuthManager{signingKey: signingKey}
}

func (a *AuthManager) MakeAuthn(userID uint) (string, error) {
	expTime := time.Now().Add(tokenLifetime)
	claims := jwt.RegisteredClaims{
		Subject:   fmt.Sprint(userID),
		ExpiresAt: jwt.NewNumericDate(expTime),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(a.signingKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *AuthManager) FetchAuthn(tknString string) (*map[string]string, error) {
	claims := jwt.RegisteredClaims{}

	tkn, err := jwt.ParseWithClaims(tknString, &claims, func(token *jwt.Token) (any, error) {
		return a.signingKey, nil
	})
	if err != nil {
		return &map[string]string{}, err
	}
	if !tkn.Valid {
		return &map[string]string{}, errors.New("invalid auth token")
	}

	return convertClaimsToMap(&claims), nil
}

func convertClaimsToMap(claims *jwt.RegisteredClaims) *map[string]string {
	sub := claims.Subject
	exp := claims.ExpiresAt
	res := map[string]string{
		"sub": sub,
		"exp": fmt.Sprint(exp),
	}
	return &res
}
