package patient

import (
	"LipidClinic/internal/lib/logger/sl"
	"LipidClinic/internal/models"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type PatientProvider2 interface {
	PatientsAll() ([]*models.Patient, error)
}

func All(log *slog.Logger, patientProvider PatientProvider2) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.patient.all"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", requestid.Get(c)),
		)

		patients, err := patientProvider.PatientsAll()
		if err != nil {
			log.Error("failed to get user", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
		}
		c.JSON(http.StatusOK, patients)

	}
}
