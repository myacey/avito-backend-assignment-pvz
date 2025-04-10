package jwt_token

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func CraeteDummyToken(role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenStr, err := token.SignedString(os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func CraeteUserToken(id uuid.UUID, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uuid": id.String(),
		"role": role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenStr, err := token.SignedString(os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func VerifyToken(tokenStr string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return os.Getenv("JWT_SECRET_KEY"), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token.Header, nil
}
