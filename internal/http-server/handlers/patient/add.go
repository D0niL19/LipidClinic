package patient

import (
	"LipidClinic/internal/lib/logger/sl"
	"LipidClinic/internal/models"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type PatientAdder interface {
	AddPatient(patient *models.Patient) error
}

func Add(log *slog.Logger, patientAdder PatientAdder) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.patient.addPatient"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", c.GetString("request_id")),
		)

		var patient *models.Patient
		if err := c.ShouldBindJSON(&patient); err != nil {
			log.Error("can not bind to json", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "Invalid inputs. Please check your inputs",
			})
			return
		}

		if err := patientAdder.AddPatient(patient); err != nil {
			log.Error("can not add patient", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
			return
		}

		log.Info("Patient successfully added")
		c.JSON(http.StatusCreated, gin.H{})

	}
}
