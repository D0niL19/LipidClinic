package auth

import (
	"LipidClinic/internal/config"
	"LipidClinic/internal/lib/logger/sl"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
)

type UserProvider interface {
	UpdatePasswordUser(id int64, password string) error
}

func ChangePassword(log *slog.Logger, userProvider UserProvider, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.auth.ChangePassword"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", requestid.Get(c)),
		)

		tokenString := c.Param("token")

		var pass struct {
			Password string `json:"password" binding:"required,min=8"`
		}

		if err := c.ShouldBindJSON(&pass); err != nil {
			log.Error("Can not bind json", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Debug("Args:", slog.String("password", pass.Password))

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.Secret), nil
		})
		if err != nil {
			log.Info("Can not parse token", sl.Err(err))
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		var id float64

		//TODO: add table for used tokens

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if id, ok = claims["uid"].(float64); !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
				return
			}

			if _, ok = claims["email"].(string); !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
				return
			}
		} else {
			log.Info("Can not parse refresh token", sl.Err(err))
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		hashPass, err := bcrypt.GenerateFromPassword([]byte(pass.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("can not generate password hash", sl.Err(err))
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		if err = userProvider.UpdatePasswordUser(int64(id), string(hashPass)); err != nil {
			log.Info("Can not change password", sl.Err(err))
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, gin.H{})

	}
}
