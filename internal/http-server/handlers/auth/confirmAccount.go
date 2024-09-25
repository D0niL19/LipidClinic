package auth

import (
	"LipidClinic/internal/lib/logger/sl"
	"LipidClinic/internal/models"
	"LipidClinic/internal/storage"
	"errors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"time"
)

type UserAdder interface {
	AddUser(user *models.User) error
	TempUserByEmail(email string) (*models.TempUser, error)
	DeleteTempUser(id int64) error
}

func ConfirmAccount(log *slog.Logger, userAdder UserAdder) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.auth.ConfirmAccount"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", requestid.Get(c)),
		)

		email := c.Param("email")
		token := c.Param("token")

		tempUser, err := userAdder.TempUserByEmail(email)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				log.Error("temp user not found", sl.Err(err))
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			log.Error("failed to get temp user", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
			return
		}

		user := models.User{
			Email:          tempUser.Email,
			HashedPassword: tempUser.HashedPassword,
			Role:           "user",
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
		}

		if tempUser.Token != token {
			log.Error("invalid token")
			log.Info(tempUser.Token)
			log.Info(token)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
			return
		}

		if err = userAdder.AddUser(&user); err != nil {
			log.Error("failed to add user", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
		}

		if err = userAdder.DeleteTempUser(tempUser.Id); err != nil {
			log.Error("failed to delete temp user", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
		}

		log.Info("successfully confirmed account")

		c.JSON(http.StatusOK, gin.H{})
	}
}
