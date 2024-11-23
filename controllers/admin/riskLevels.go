package admin

import (
	"almanac-api/collections"
	"almanac-api/db"
	"almanac-api/middleware"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/almanac-app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CreateRiskLevel struct {
	Level          int      `json:"level"`
	Municipalities []string `json:"municipalities"`
}

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

func (r RiskLevelsController) Create(c *gin.Context) {
	mongoClient := db.GetDbClient()

	user, ok := middleware.GetUserFromContext(c)
	if !ok || !user.IsAdmin() {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Unauthorized": "Not admin"})
		return
	}

	var createRiskLevel CreateRiskLevel
	err := c.BindJSON(&createRiskLevel)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid request body"})
		return
	}

	municipalityObjectIDs := make([]primitive.ObjectID, len(createRiskLevel.Municipalities))
	for i, idStr := range createRiskLevel.Municipalities {
		objectID, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid municipality ID format: %s", idStr)})
			return
		}
		municipalityObjectIDs[i] = objectID
	}

	var municipalities []models.Municipality
	cursor, err := collections.GetMunicipalityCollection(*mongoClient).Find(c, bson.M{
		"_id": bson.M{"$in": municipalityObjectIDs},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch municipalities"})
		return
	}
	defer cursor.Close(c)

	if err = cursor.All(c, &municipalities); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode municipalities"})
		return
	}

	if len(municipalities) != len(createRiskLevel.Municipalities) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some municipality IDs are invalid"})
		return
	}

	riskLevel := models.RiskLevel{
		User:           user.Id,
		Municipalities: municipalities,
		Level:          createRiskLevel.Level,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	result, err := collections.GetRiskLevelCollection(*mongoClient).InsertOne(c, riskLevel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	var createdRiskLevel models.RiskLevel
	err = collections.GetRiskLevelCollection(*mongoClient).FindOne(c, bson.D{{"_id", result.InsertedID}}).Decode(&createdRiskLevel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, createdRiskLevel)
}

func (r RiskLevelsController) Delete(c *gin.Context) {
	mongoClient := db.GetDbClient()
	id := c.Param("id")
	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "no valid id provided"})
		return
	}

	update := bson.M{
		"updatedAt":  time.Now(),
		"archivedAt": time.Now(),
	}
	result, err := collections.GetRiskLevelCollection(*mongoClient).UpdateOne(c, bson.M{"_id": objId}, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if result.MatchedCount != 1 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "error updating risk level"})
		return
	}

	c.JSON(http.StatusOK, bson.M{})
}
