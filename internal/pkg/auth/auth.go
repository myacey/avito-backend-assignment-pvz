package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/response"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/entity"
)

type TokenChecker interface {
	VerifyToken(token string) (map[string]interface{}, error)
}

type AuthService struct {
	tokenSrv TokenChecker
}

func New(tokenSrv TokenChecker) *AuthService {
	return &AuthService{
		tokenSrv: tokenSrv,
	}
}

func getToken(ctx *gin.Context) (string, error) {
	bearerToken := ctx.GetHeader("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) != 2 {
		return "", errors.New("invalid token")
	}

	tokenStr := splitToken[1]

	return tokenStr, nil
}

func (s *AuthService) AuthMiddleware(neededRole entity.Role) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := getToken(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Error{
				Code:      http.StatusUnauthorized,
				Message:   err.Error(),
				RequestId: ctx.GetHeader("X-Request-Id"),
			})
			return
		}

		claims, err := s.tokenSrv.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Error{
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
