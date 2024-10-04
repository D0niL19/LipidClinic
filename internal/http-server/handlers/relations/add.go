package relations

import (
	"LipidClinic/internal/lib/logger/sl"
	"LipidClinic/internal/models"
	"LipidClinic/internal/storage"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type RelationAdder interface {
	AddRelation(relation *models.Relation) error
}

func Add(log *slog.Logger, relationtAdder RelationAdder) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.relations.add"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", c.GetString("request_id")),
		)

		var relation *models.Relation
		if err := c.ShouldBindJSON(&relation); err != nil {
			log.Error("can not bind to json", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "Invalid inputs. Please check your inputs",
			})
			return
		}

		if err := relationtAdder.AddRelation(relation); err != nil {
			if errors.Is(err, storage.ErrRelationExists) {
				log.Error("relation already exists", sl.Err(err))
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{})
				return
			}
			log.Error("can not add patient", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
			return
		}

		log.Info("Patient successfully added")
		c.JSON(http.StatusCreated, gin.H{})

	}
}
