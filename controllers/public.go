package controllers

import (
	"almanac-api/collections"
	"almanac-api/db"
	"almanac-api/serializers"
	"almanac-api/services"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gitlab.com/almanac-app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	case "7days":
		sevenDaysAgo := primitive.NewDateTimeFromTime(time.Now().AddDate(0, 0, -7))
		filter["timestamp"] = bson.M{"$gte": sevenDaysAgo}
	case "30days":
		thirtyDaysAgo := primitive.NewDateTimeFromTime(time.Now().AddDate(0, 0, -30))
		filter["timestamp"] = bson.M{"$gte": thirtyDaysAgo}
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
	// limit to 2 pins
	projection := bson.D{{"geoData", bson.D{{"features", bson.D{{"$slice", 2}}}}}}
	opts := options.Find().SetSort(bson.D{{"timestamp", -1}}).SetProjection(projection)
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

func (p PublicController) Pois(c *gin.Context) {
	mongoClient := db.GetDbClient()
	filter := bson.D{{"active", true}}
	cur, err := collections.GetPoicollection(*mongoClient).Find(c, filter)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pois := make([]models.Poi, 0)
	err = cur.All(c, &pois)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pois)
}

func (p PublicController) RiskLevels(c *gin.Context) {
	riskLevelService := services.RiskLevelService{C: c}

	response, err := riskLevelService.PublicRiskLevels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (p PublicController) Report(c *gin.Context) {
	mongoClient := db.GetDbClient()
	dateString := c.Query("date")
	filter := bson.M{
		"archivedAt": nil,
		"date":       dateString,
	}
	opts := options.FindOne().SetSort(bson.D{{"date", -1}})

	var report models.DailyReport
	err := collections.GetDailyReportsCollection(*mongoClient).FindOne(c, filter, opts).Decode(&report)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No report for this date"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	serializer := serializers.ReportSerializer{
		DailyReport: report,
	}
	c.JSON(http.StatusOK, serializer.PublicResponse())
}
