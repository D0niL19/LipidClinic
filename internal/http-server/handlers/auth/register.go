package auth

import (
	"LipidClinic/internal/config"
	"LipidClinic/internal/lib/email"
	"LipidClinic/internal/lib/logger/sl"
	"LipidClinic/internal/models"
	"LipidClinic/internal/storage"
	"errors"
	"fmt"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
	"time"
)

type TempUserAdder interface {
	AddTempUser(tempUser *models.TempUser) error
	UserByEmail(email string) (*models.User, error)
	UpdateTempUser(tempUser *models.TempUser) error
}

func Register(log *slog.Logger, tempUserAdder TempUserAdder, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.auth.Register"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", requestid.Get(c)),
		)

		var input models.SignInUser

		if err := c.ShouldBindJSON(&input); err != nil {
			log.Error("can not bind to json", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "Invalid inputs. Please check your inputs",
			})
			return
		}

		user, err := tempUserAdder.UserByEmail(input.Email)
		if user != nil {
			log.Error("user already exists")
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "email already exists"})
			return
		}

		hashPass, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("can not generate password hash", sl.Err(err))
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		token, err := email.GenerateRandomAuthString()
		if err != nil {
			log.Error("can not generate auth token", sl.Err(err))
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		tempUser := models.TempUser{
			Email:          input.Email,
			HashedPassword: string(hashPass),
			Token:          token,
			CreatedAt:      time.Now().UTC(),
		}

		err = tempUserAdder.UpdateTempUser(&tempUser)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				// Попытка добавить пользователя, если не найден
				if err = tempUserAdder.AddTempUser(&tempUser); err != nil {
					log.Error("can not add temp user", sl.Err(err))
					c.JSON(http.StatusInternalServerError, gin.H{})
					return
				}
			} else {
				// Ошибка при обновлении пользователя
				log.Error("can not update temp user", sl.Err(err))
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
		}

		url := fmt.Sprintf("localhost:8080/auth/confirm-account/%s/%s", tempUser.Email, token)

		message := email.EmailConfirmMessage(url)

		if err = email.SendEmail(input.Email, message, cfg); err != nil {
			log.Error("can not send email", sl.Err(err))
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		log.Info("email sent")

		c.JSON(http.StatusOK, gin.H{})
	}
}
