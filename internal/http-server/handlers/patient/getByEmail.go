package patient

import (
	"LipidClinic/internal/lib/logger/sl"
	"LipidClinic/internal/models"
	"LipidClinic/internal/storage"
	"errors"

	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type PatientProvider interface {
	PatientByEmail(email string) (*models.Patient, error)
}

func ByEmail(log *slog.Logger, provider PatientProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.patient.getByEmail"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", c.GetString("request_id")),
		)

		email := c.Query("email")
		if email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
			return
		}

		patient, err := provider.PatientByEmail(email)
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
