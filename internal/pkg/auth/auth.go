package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/myacey/avito-backend-assignment-pvz/internal/httpserver/handler"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/response"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/entity"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/jwttoken"
)

const (
	HeaderAuthorization = "Authorization"

	CtxKeyUserType = "User-Type"
)

type TokenChecker interface {
	VerifyToken(token string) (map[string]interface{}, error)
}

type Service struct {
	tokenSrv TokenChecker
}

func New(tokenSrv TokenChecker) *Service {
	return &Service{
		tokenSrv: tokenSrv,
	}
}

func getToken(ctx *gin.Context) (string, error) {
	bearerToken := ctx.GetHeader(HeaderAuthorization)
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) != 2 {
		return "", errors.New("invalid token")
	}

	tokenStr := splitToken[1]

	return tokenStr, nil
}

func (s *Service) AuthMiddleware(neededRole ...entity.Role) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := getToken(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Error{
				Code:      http.StatusUnauthorized,
				Message:   err.Error(),
				RequestID: ctx.GetHeader(handler.HeaderRequestID),
			})
			return
		}

		claims, err := s.tokenSrv.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Error{
				Code:      http.StatusUnauthorized,
				Message:   err.Error(),
				RequestID: ctx.GetHeader(handler.HeaderRequestID),
			})
			return
		}

		r, ok := claims[jwttoken.JwtClaimRole]
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Error{
				Code:      http.StatusUnauthorized,
				Message:   "invalid role",
				RequestID: ctx.GetHeader(handler.HeaderRequestID),
			})
		}

		for _, nRole := range neededRole {
			if r.(string) != string(nRole) {
				continue
			}

			ctx.Set(CtxKeyUserType, claims[jwttoken.JwtClaimRole])
			ctx.Next()
			return
		}

		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Error{
			Code:      http.StatusUnauthorized,
			Message:   "invalid role",
			RequestID: ctx.GetHeader(handler.HeaderRequestID),
		})
	}
}
