package admin

import (
	"almanac-api/collections"
	"almanac-api/db"
	"almanac-api/serializers"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/almanac-app/models"
	"go.mongodb.org/mongo-driver/bson"
)

type MunicipalitiesController struct {
}

func (m MunicipalitiesController) List(c *gin.Context) {
	mongoClient := db.GetDbClient()
	cur, err := collections.GetMunicipalityCollection(*mongoClient).Find(c, bson.D{{}})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	municipalities := make([]models.Municipality, 0)
	err = cur.All(c, &municipalities)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	serializer := serializers.MunicipalitiesSerializer{
		Municipalities: municipalities,
	}
	c.JSON(http.StatusOK, serializer.Response())
}
