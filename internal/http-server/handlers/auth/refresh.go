package auth

import (
	"LipidClinic/internal/config"
	"LipidClinic/internal/lib/logger/sl"
	jwtlib "LipidClinic/internal/lib/tokens/jwt"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"net/http"
)

func Refresh(log *slog.Logger, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.auth.Refresh"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", requestid.Get(c)),
		)

		refreshTokenString, err := c.Cookie("refresh_token")
		if err != nil {
			log.Info("Can not get refresh token", sl.Err(err))
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		if refreshTokenString == "" {
			log.Info("Can not get refresh token", sl.Err(err))
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.Secret), nil
		})
		if err != nil {
			log.Info("Can not parse refresh token", sl.Err(err))
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		var id float64
		var email string

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if id, ok = claims["uid"].(float64); !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
				return
			}

			if email, ok = claims["email"].(string); !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
				return
			}
		} else {
			log.Info("Can not parse refresh token", sl.Err(err))
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		accessToken, err := jwtlib.NewToken(int64(id), email, cfg.Jwt.Access, cfg.Jwt.Secret)
		if err != nil {
			log.Error("failed to generate token", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": accessToken,
		})
	}
}
