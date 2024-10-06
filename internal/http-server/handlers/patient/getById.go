package patient

import (
	"LipidClinic/internal/lib/logger/sl"
	"LipidClinic/internal/models"
	"LipidClinic/internal/storage"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
)

type PatientProvider1 interface {
	PatientById(id int64) (*models.Patient, error)
}

func ById(log *slog.Logger, provider PatientProvider1) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.patient.getByEmail"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", c.GetString("request_id")),
		)

		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
			return
		}

		atoi, err := strconv.Atoi(id)
		if err != nil {
			log.Error("Can not parse id")
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		patient, err := provider.PatientById(int64(atoi))
		if err != nil {
			if errors.Is(err, storage.ErrPatientNotFound) {
				log.Info("Patient not found")
				c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
			}
			log.Error("can not get patient by email", sl.Err(err))
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, patient)

	}
}
