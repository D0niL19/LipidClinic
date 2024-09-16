package add

import (
	"LipidClinic/internal/lib/logger/sl"
	"LipidClinic/internal/models"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@v2.45.1 --name=UserAdder
type UserAdder interface {
	AddUser(user *models.User) error
}

func New(log *slog.Logger, userAdder UserAdder) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.user.add.New"

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
			log.Error("failed to add user", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
			return
		}

		log.Info("User successfully added")
		c.JSON(http.StatusCreated, gin.H{"status": "ok"})
	}
}
