package admin

import (
	"almanac-api/collections"
	"almanac-api/db"
	"almanac-api/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/almanac-app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NewsItemsController struct {
}

func (n NewsItemsController) List(c *gin.Context) {
	mongoClient := db.GetDbClient()
	opts := options.Find().SetSort(bson.D{{"createdAt", -1}})
	cur, err := collections.GetNewsItemCollection(*mongoClient).Find(c, bson.D{{}}, opts)

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

func (n NewsItemsController) Create(c *gin.Context) {
	mongoClient := db.GetDbClient()

	user, ok := middleware.GetUserFromContext(c)
	if !ok || !user.IsAdmin() {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Unauthorized": "Not admin"})
		return
	}

	var newsItem models.NewsItem
	err := c.BindJSON(&newsItem)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid request body"})
		return
	}

	newsItem.CreatedAt = time.Now()
	newsItem.UpdatedAt = time.Now()
	newsItem.User = user.Id
	newsItem.Provider = "manual"

	result, err := collections.GetNewsItemCollection(*mongoClient).InsertOne(c, newsItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	var createdNewsItem models.NewsItem
	err = collections.GetNewsItemCollection(*mongoClient).FindOne(c, bson.D{{"_id", result.InsertedID}}).Decode(&createdNewsItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, createdNewsItem)
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

	update := bson.M{
		"title":       updateNewsItem.Title,
		"description": updateNewsItem.Description,
		"status":      updateNewsItem.Status,
		"category":    updateNewsItem.Category,
		"geoData":     updateNewsItem.GeoData,
		"timestamp":   updateNewsItem.Timestamp,
		"updatedAt":   time.Now(),
	}
	result, err := collections.GetNewsItemCollection(*mongoClient).UpdateOne(c, bson.M{"_id": objId}, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if result.MatchedCount != 1 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "error updating news item"})
		return
	}

	var newsItem models.NewsItem
	err = collections.GetNewsItemCollection(*mongoClient).FindOne(c, bson.D{{"_id", objId}}).Decode(&newsItem)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "newsitem not found"})
		return
	}

	c.JSON(http.StatusOK, newsItem)
}
