package patient

import (
	"LipidClinic/internal/lib/logger/sl"
	"LipidClinic/internal/storage"
	"errors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
)

type PatientDeleter interface {
	DeletePatient(id int) error
}

func Delete(log *slog.Logger, patientDeleter PatientDeleter) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.patient.delete"

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

		err := patientDeleter.DeletePatient(idInt)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				log.Info("user not found")
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			log.Error("failed to get user", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
		}
		log.Info("user deleted", slog.String("id", id))
		c.JSON(http.StatusOK, gin.H{"status": "ok"})

	}
}
