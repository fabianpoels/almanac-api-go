package admin

import (
	"almanac-api/db"
	"almanac-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NewsItemsController struct {
}

func (n NewsItemsController) List(c *gin.Context) {
	mongoClient := db.GetDbClient()
	opts := options.Find().SetSort(bson.D{{"createdAt", -1}})
	cur, err := models.GetNewsItemCollection(*mongoClient).Find(c, bson.D{{}}, opts)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newsItems := make([]models.NewsItem, 0)
	err = cur.All(c, &newsItems)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newsItems)
}

// func (n NewsItemsController) Patch(c *gin.Context) {

// }
