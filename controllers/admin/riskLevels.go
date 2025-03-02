package admin

import (
	"almanac-api/collections"
	"almanac-api/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/almanac-app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RiskLevelsController struct {
}

func (r RiskLevelsController) List(c *gin.Context) {
	mongoClient := db.GetDbClient()
	filter := bson.M{"archivedAt": nil}
	opts := options.Find().SetSort(bson.D{{"createdAt", -1}})
	cur, err := collections.GetRiskLevelCollection(*mongoClient).Find(c, filter, opts)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	riskLevels := make([]models.RiskLevel, 0)
	err = cur.All(c, &riskLevels)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, riskLevels)
}
