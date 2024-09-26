package controllers

import (
	"almanac-api/collections"
	"almanac-api/db"
	"almanac-api/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gitlab.com/almanac-app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PublicController struct {
}

type newsItemRequest struct {
	Span string `json:"span"`
	From string `json:"from"`
	To   string `json:"to"`
}

func (p PublicController) NewsItems(c *gin.Context) {
	var req newsItemRequest
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid request body"})
		return
	}

	filter := bson.M{"status": "published"}
	switch req.Span {
	case "12hr":
		twelveHoursAgo := primitive.NewDateTimeFromTime(time.Now().Add(-12 * time.Hour))
		filter["timestamp"] = bson.M{"$gte": twelveHoursAgo}
	case "24hr":
		oneDayAgo := primitive.NewDateTimeFromTime(time.Now().Add(-24 * time.Hour))
		filter["timestamp"] = bson.M{"$gte": oneDayAgo}
	case "48hr":
		twoDaysAgo := primitive.NewDateTimeFromTime(time.Now().Add(-48 * time.Hour))
		filter["timestamp"] = bson.M{"$gte": twoDaysAgo}
	case "week":
		startOfWeek := primitive.NewDateTimeFromTime(utils.GetStartOfWeek())
		filter["timestamp"] = bson.M{"$gte": startOfWeek}
	case "month":
		startOfMonth := primitive.NewDateTimeFromTime(utils.GetStartOfMonth())
		filter["timestamp"] = bson.M{"$gte": startOfMonth}
	case "custom":
		fromDate, err := time.Parse("2006/01/02", req.From)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid request body"})
			return
		}

		toDate, err := time.Parse("2006/01/02", req.To)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid request body"})
			return
		}

		filter["timestamp"] = bson.M{
			"$gte": primitive.NewDateTimeFromTime(fromDate),
			"$lte": primitive.NewDateTimeFromTime(toDate),
		}
	}

	mongoClient := db.GetDbClient()
	opts := options.Find().SetSort(bson.D{{"timestamp", -1}})
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
