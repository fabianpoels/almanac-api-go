package admin

import (
	"almanac-api/db"
	"almanac-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (n NewsItemsController) Update(c *gin.Context) {
	mongoClient := db.GetDbClient()
	id := c.Param("id")
	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "no valid id provided"})
		return
	}

	var updateNewsItem models.NewsItem
	err = c.BindJSON(&updateNewsItem)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid request body"})
		return
	}

	update := bson.M{"status": updateNewsItem.Status, "geoData": updateNewsItem.GeoData, "updatedAt": time.Now()}
	result, err := models.GetNewsItemCollection(*mongoClient).UpdateOne(c, bson.M{"_id": objId}, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if result.MatchedCount != 1 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "error updating news item"})
		return
	}

	var newsItem models.NewsItem
	err = models.GetNewsItemCollection(*mongoClient).FindOne(c, bson.D{{"_id", objId}}).Decode(&newsItem)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "newsitem not found"})
		return
	}

	c.JSON(http.StatusOK, newsItem)
}
