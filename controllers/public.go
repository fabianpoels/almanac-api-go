package controllers

import (
	"almanac-api/db"
	"almanac-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type PublicController struct {
}

func (p PublicController) NewsItems(c *gin.Context) {
	mongoClient := db.GetDbClient()
	cur, err := models.GetNewsItemCollection(*mongoClient).Find(c, bson.M{})

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
	cur, err := models.GetCategoryCollection(*mongoClient).Find(c, bson.M{})

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
