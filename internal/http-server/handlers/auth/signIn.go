package auth

import (
	"LipidClinic/internal/config"
	"LipidClinic/internal/lib/logger/sl"
	"LipidClinic/internal/lib/tokens/jwt"
	"LipidClinic/internal/models"
	"LipidClinic/internal/storage"
	"errors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
)

type UserGetter interface {
	UserByEmail(email string) (*models.User, error)
}

func SignIn(log *slog.Logger, userGetter UserGetter, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.auth.SignIn"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", requestid.Get(c)),
		)

		var user *models.SignInUser
		if err := c.ShouldBindJSON(&user); err != nil {
			log.Error("can not bind to json", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "Invalid inputs. Please check your inputs",
			})
			return
		}

		userdb, err := userGetter.UserByEmail(user.Email)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				log.Info("user not found")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
				return
			}
			log.Error("failed to get user", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
			return
		}

		if err = bcrypt.CompareHashAndPassword([]byte(userdb.HashedPassword), []byte(user.Password)); err != nil {
			log.Info("wrong password", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
			return
		}

		accessToken, err := jwt.NewToken(userdb.Id, userdb.Email, cfg.Jwt.Access, cfg.Jwt.Secret)
		if err != nil {
			log.Error("failed to generate token", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
			return
		}
		refreshToken, err := jwt.NewToken(userdb.Id, userdb.Email, cfg.Jwt.Refresh, cfg.Jwt.Secret)
		if err != nil {
			log.Error("failed to generate token", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.SetCookie(
			"refresh_token",
			refreshToken,
			int(cfg.Jwt.Refresh.Seconds()),
			"/",
			cfg.HttpServer.Address,
			false,
			true,
		)

		c.JSON(http.StatusOK, gin.H{
			"token": accessToken,
		})
		log.Info("User logged in")
	}
}
