package relations

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

type RelationProvider interface {
	Relation(id int64) (*models.Relation, error)
}

func ById(log *slog.Logger, provider RelationProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.relations.Get"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", c.GetString("request_id")),
		)

		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}
		intId, err := strconv.Atoi(id)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		relation, err := provider.Relation(int64(intId))
		if err != nil {
			if errors.Is(err, storage.ErrRelationNotFound) {
				log.Info("Relation not found")
				c.JSON(http.StatusNotFound, gin.H{"error": "Relation not found"})
			}
			log.Error("can not get relation by id", sl.Err(err))
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, relation)

	}
}
