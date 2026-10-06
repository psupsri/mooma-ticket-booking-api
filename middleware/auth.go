package middleware

import (
	"fmt"
	"mooma-api/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(ctx *gin.Context) {
	key := config.Load().JWTSecret
	authHeader := ctx.GetHeader("Authorization")
	accessToken := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(key), nil
	})
	if err != nil || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "failed",
			"message": "unauthorized",
		})
		ctx.Abort()
	}
	ctx.Next()
}
