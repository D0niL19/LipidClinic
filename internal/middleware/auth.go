package middleware

import (
	"LipidClinic/internal/lib/logger/sl"
	"errors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"net/http"
	"strings"
)

func AuthMiddleware(jwtKey string, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "middleware.AuthMiddleware"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", requestid.Get(c)),
		)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Error("Empty authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Error("Invalid authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})

		switch {
		case token.Valid:
			log.Info("User authenticated successfully")
			c.Next()
		case errors.Is(err, jwt.ErrTokenMalformed):
			log.Error("Token is malformed", sl.Err(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is malformed"})
			c.Abort()
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			log.Error("Token signature is invalid", sl.Err(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token signature is invalid"})
			c.Abort()
		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet): // Token is either expired or not active yet
			log.Error("Token is not yet valid", sl.Err(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is not yet valid"})
			c.Abort()
		default:
			log.Error("Invalid token", sl.Err(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
		}

	}
}
