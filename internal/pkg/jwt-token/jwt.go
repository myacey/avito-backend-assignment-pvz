package jwt_token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenServiceConfig struct {
	SecretKey string `mapstructure:"jwt_secret_key"`
}

type JWTTokenService struct {
	secretKey []byte
}

func New(cfg TokenServiceConfig) *JWTTokenService {
	return &JWTTokenService{[]byte(cfg.SecretKey)}
}

func (s *JWTTokenService) CraeteDummyToken(role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenStr, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (s *JWTTokenService) CraeteUserToken(id uuid.UUID, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uuid": id.String(),
		"role": role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenStr, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (s *JWTTokenService) VerifyToken(tokenStr string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
