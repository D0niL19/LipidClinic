package patient

import (
	"LipidClinic/internal/lib/logger/sl"
	"LipidClinic/internal/models"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type PatientProvider interface {
	PatientByEmail(email string) (*models.Patient, error)
}

func GetByEmail(log *slog.Logger, provider PatientProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.patient.getByEmail"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", c.GetString("request_id")),
		)

		email := c.Param("email")

		patient, err := provider.PatientByEmail(email)
		if err != nil {
			log.Error("can not get patient by email", sl.Err(err))
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, patient)

	}
}
