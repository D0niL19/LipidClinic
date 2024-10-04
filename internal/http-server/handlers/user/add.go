package user

import (
	"LipidClinic/internal/lib/logger/sl"
	"LipidClinic/internal/models"
	"LipidClinic/internal/storage"
	"errors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type UserAdder interface {
	AddUser(user *models.User) error
}

func Add(log *slog.Logger, userAdder UserAdder) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.user.add"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", requestid.Get(c)),
		)

		var user *models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			log.Error("can not bind to json", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "Invalid inputs. Please check your inputs",
			})
			return
		}

		if err := userAdder.AddUser(user); err != nil {
			if errors.Is(err, storage.ErrUserExists) {
				log.Error("user already exists", sl.Err(err))
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{})
				return
			}
			log.Error("failed to add user", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
			return
		}

		log.Info("User successfully added")
		c.JSON(http.StatusCreated, gin.H{"status": "ok"})
	}
}
