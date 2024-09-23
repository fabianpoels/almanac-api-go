package controllers

import (
	"almanac-api/collections"
	"almanac-api/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/almanac-app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PublicController struct {
}

func (p PublicController) NewsItems(c *gin.Context) {
	mongoClient := db.GetDbClient()
	twoDaysAgo := primitive.NewDateTimeFromTime(time.Now().AddDate(0, 0, -2))
	filter := bson.M{
		"createdAt": bson.M{"$gte": twoDaysAgo},
		"status":    "published",
	}

	opts := options.Find().SetSort(bson.D{{"createdAt", -1}})
	cur, err := collections.GetNewsItemCollection(*mongoClient).Find(c, filter, opts)

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

func (p PublicController) Categories(c *gin.Context) {
	mongoClient := db.GetDbClient()
	filter := bson.D{{"active", true}}
	cur, err := collections.GetCategoryCollection(*mongoClient).Find(c, filter)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	categories := make([]models.Category, 0)
	err = cur.All(c, &categories)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}
