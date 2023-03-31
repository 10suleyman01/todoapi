package util

import (
	"crypto/sha1"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"time"
)

// 6L6v6S6o5H
// aB2pC0pZ6ceT

const (
	ApiV1 = "/api"
	ApiV2 = "/api/v2"
)

const (
	Postgres = "postgres"
	Server   = "server"
	Token    = "token"
)

func ConfigPath(path, key string) string {
	return fmt.Sprintf("%s.%s", path, key)
}

func GeneratePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(viper.GetString("secret.salt"))))
}

func GenerateToken(ttl time.Duration, payload interface{}, secretJWTKey string) (string, error) {

	jwtGen := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := jwtGen.Claims.(jwt.MapClaims)

	claims["sub"] = payload
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwtGen.SignedString([]byte(secretJWTKey))
	if err != nil {
		return "", fmt.Errorf("generating JWT token failed: %w", err)
	}

	return token, nil
}

func ValidateToken(token, secretJwtKey string) (interface{}, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method: %s", token.Header["alg"])
		}
		return []byte(secretJwtKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalidate token: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token claim")
	}

	return claims["sub"], nil
}
