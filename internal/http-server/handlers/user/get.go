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
	"strconv"
)

type UserGetter interface {
	GetUser(id int) (*models.User, error)
}

func Get(log *slog.Logger, userGetter UserGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.user.get.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", requestid.Get(c)),
		)

		id := c.Param("id")

		idInt, err2 := strconv.Atoi(id)
		if err2 != nil {
			log.Info("id is not a number")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id is not a number"})
			return
		}

		user, err := userGetter.GetUser(idInt)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				log.Info("user not found")
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			log.Error("failed to get user", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
			return
		}
		log.Info("User successfully got")
		c.JSON(http.StatusOK, user)
	}
}
