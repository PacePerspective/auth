package bajwt

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/paceperspective/googlesecret"
)

var (
	StandardTokenLife = time.Hour * 1
	ProjectID         string
)

func Create(ctx context.Context, userName string, expiryDuration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": userName,
			"exp":      time.Now().Add(expiryDuration).Unix(),
		})

	k, err := getSecretTokenKey(ctx)
	if err != nil {
		return "", err
	}
	tokenString, err := token.SignedString(k)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetTokenFromHttpHeader(header string) (string, error) {
	split := strings.Split(header, "Bearer ")
	if len(split) != 2 {
		return "", errors.New("invalid token format")
	}
	return split[1], nil
}

func Verify(ctx context.Context, tokenString string) error {
	token, err := jwt.Parse(tokenString, func(_ *jwt.Token) (interface{}, error) {
		return getSecretTokenKey(ctx)
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("token invalid")
	}
	return nil
}

func GetStringClaimFromToken(ctx context.Context, tokenString, key string) (string, error) {
	token, err := jwt.Parse(tokenString, func(_ *jwt.Token) (interface{}, error) {
		return getSecretTokenKey(ctx)
	})
	if err != nil {
		return "", err
	}
	claims := token.Claims.(jwt.MapClaims)
	return claims[key].(string), nil
}

func getSecretTokenKey(ctx context.Context) ([]byte, error) {
	secret, err := googlesecret.New(ctx, ProjectID, "jwt-auth-token-key", "latest")
	if err != nil {
		return nil, fmt.Errorf("failed to get secret key: %w", err)
	}
	return []byte(secret.Value), nil
}
