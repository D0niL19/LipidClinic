package auth

import (
	"LipidClinic/internal/lib/logger/sl"
	"LipidClinic/internal/models"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type UserProvider interface {
	UserByEmail(email string) (*models.User, error)
}

func ForgetPassword(log *slog.Logger, userProvider UserProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.auth.ForgetPassword"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", c.GetString("request_id")),
		)

		var email string

		if err := c.ShouldBindJSON(email); err != nil {
			log.Error("can not bind to json", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "Invalid inputs. Please check your inputs",
			})
			return
		}

		_, err := userProvider.UserByEmail(email)
		if err != nil {
			log.Error("can not find user by email", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{})
			return
		}

	}
}
