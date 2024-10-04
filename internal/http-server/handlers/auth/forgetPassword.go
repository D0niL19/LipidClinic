package auth

import (
	"LipidClinic/internal/config"
	"LipidClinic/internal/lib/email"
	"LipidClinic/internal/lib/logger/sl"
	"LipidClinic/internal/lib/tokens/jwt"
	"LipidClinic/internal/models"
	"fmt"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type UserChanger interface {
	UserByEmail(email string) (*models.User, error)
}

func ForgetPassword(log *slog.Logger, userProvider UserChanger, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.auth.ForgetPassword"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", requestid.Get(c)),
		)

		var json struct {
			Email string `json:"email" binding:"required"`
		}

		if err := c.ShouldBindJSON(&json); err != nil {
			log.Error("can not bind to json", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "Invalid inputs. Please check your inputs",
			})
			return
		}

		userEmail := json.Email

		user, err := userProvider.UserByEmail(userEmail)
		if err != nil {
			log.Error("can not find user by email", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{})
			return
		}

		tokenString, err := jwt.NewToken(user.Id, user.Email, cfg.TokenExpiration, cfg.Jwt.Secret)
		if err != nil {
			log.Error("can not create token", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
		}
		url := fmt.Sprintf("localhost:8080/auth/change-password/%s", tokenString)

		message := email.PasswordChangeMessage(url)

		if err = email.SendEmail(userEmail, message, cfg); err != nil {
			log.Error("can not send email", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
			return
		}
		log.Info("reset password email sent")

		c.Status(http.StatusOK)
	}
}
