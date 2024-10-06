package relations

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

type RelationDeleter interface {
	DeleteRelation(patId, relId int64) error
}

func Delete(log *slog.Logger, relationDeleter RelationDeleter) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.relation.delete"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", requestid.Get(c)),
		)

		patId := c.Query("patient_id")
		relId := c.Query("relative_id")

		patIdint, err := strconv.Atoi(patId)
		if err != nil {
			log.Info("id is not a number")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "patient_id is not a number"})
			return
		}

		relIdint, err := strconv.Atoi(relId)
		if err != nil {
			log.Info("id is not a number")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "relative_id is not a number"})
			return
		}

		err = relationDeleter.DeleteRelation(int64(patIdint), int64(relIdint))
		if err != nil {
			if errors.Is(err, storage.ErrRelationNotFound) {
				log.Info("relation not found")
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "relation not found"})
				return
			}
			log.Error("failed to delete relation", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
		}
		//log.Info("user deleted", slog.String("id", id))
		c.JSON(http.StatusOK, gin.H{"status": "ok"})

	}
}
