package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/response"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/entity"
	jwt_token "github.com/myacey/avito-backend-assignment-pvz/internal/pkg/jwt-token"
)

// type TokenService interface {
// 	CreateToken(role string, expireTime time.Time) (string, error)
// 	VerifyToken(token string) (map[string]interface{}, error)
// }

// type AuthService struct {
// 	tokenSrv TokenService
// }

// func NewAuthService(tokenSrv TokenService) *AuthService {
// 	return
// }

func getToken(ctx *gin.Context) (string, error) {
	bearerToken := ctx.GetHeader("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) != 2 {
		return "", errors.New("invalid token")
	}

	tokenStr := splitToken[1]

	return tokenStr, nil
}

func AuthMiddleware(neededRole entity.Role) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := getToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, response.Error{
				Code:      http.StatusUnauthorized,
				Message:   err.Error(),
				RequestId: ctx.GetHeader("X-Request-Id"),
			})
			return
		}

		claims, err := jwt_token.VerifyToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, response.Error{
				Code:      http.StatusUnauthorized,
				Message:   err.Error(),
				RequestId: ctx.GetHeader("X-Request-Id"),
			})
			return
		}

		if claims["role"].(string) == string(neededRole) {
			ctx.Set("User-Type", claims["role"])
			ctx.Next()
			return
		}

		ctx.JSON(http.StatusUnauthorized, response.Error{
			Code:      http.StatusUnauthorized,
			Message:   "invalid role",
			RequestId: ctx.GetHeader("X-Request-Id"),
		})
	}
}
