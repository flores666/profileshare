package security

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type JWTService struct {
	accessSecret  []byte
	refreshSecret []byte
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

func NewJWTService(settings Settings) *JWTService {
	return &JWTService{
		accessSecret:  []byte(settings.AccessSecret),
		refreshSecret: []byte(settings.RefreshSecret),
		AccessTTL:     time.Duration(settings.AccessTTL) * time.Minute,
		RefreshTTL:    time.Duration(settings.RefreshTTL) * 24 * time.Hour,
	}
}

func (s *JWTService) GenerateTokens(userId string) (*TokenPair, error) {
	accessClaims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(s.AccessTTL).Unix(),
		"type":    "access",
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(s.accessSecret)
	if err != nil {
		return nil, err
	}

	refreshToken := generateSecureToken(32)

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *JWTService) GetValue(tokenStr, key string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.accessSecret, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "access" {
		return "", errors.New("invalid claims")
	}

	value, ok := claims[key].(string)
	if !ok {
		return "", errors.New(fmt.Sprintf("invalid %s", key))
	}

	return value, nil
}

func generateSecureToken(length int) string {
	bytes := make([]byte, length)
	_, _ = rand.Read(bytes)
	return base64.RawURLEncoding.EncodeToString(bytes)
}
