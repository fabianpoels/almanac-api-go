package services

import (
	"almanac-api/collections"
	"almanac-api/db"
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/almanac-app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const cacheKey = "riskLevels:public"

type CreateRiskLevel struct {
	Level          int      `json:"level"`
	Municipalities []string `json:"municipalities"`
}

type RiskLevelService struct {
	C *gin.Context
}

func (service *RiskLevelService) Create(createRiskLevel *CreateRiskLevel, user *models.User) (riskLevel *models.RiskLevel, err error) {
	mongoClient := db.GetDbClient()

	municipalityObjectIDs := make([]primitive.ObjectID, len(createRiskLevel.Municipalities))
	for i, idStr := range createRiskLevel.Municipalities {
		objectID, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			log.Printf("Invalid municipality ID format: %s", idStr)
			return riskLevel, err
		}
		municipalityObjectIDs[i] = objectID
	}

	var municipalities []models.Municipality
	cursor, err := collections.GetMunicipalityCollection(*mongoClient).Find(service.C, bson.M{
		"_id": bson.M{"$in": municipalityObjectIDs},
	})
	if err != nil {
		log.Println("error: Failed to fetch municipalities")
		return riskLevel, err
	}
	defer cursor.Close(service.C)

	if err = cursor.All(service.C, &municipalities); err != nil {
		log.Println("error: Failed to decode municipalities")
		return riskLevel, err
	}

	if len(municipalities) != len(createRiskLevel.Municipalities) {
		log.Println("error: Some municipality IDs are invalid")
		return riskLevel, err
	}

	riskLevel = &models.RiskLevel{
		User:           user.Id,
		Municipalities: municipalities,
		Level:          createRiskLevel.Level,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	result, err := collections.GetRiskLevelCollection(*mongoClient).InsertOne(service.C, riskLevel)
	if err != nil {
		return riskLevel, err
	}

	var createdRiskLevel models.RiskLevel
	err = collections.GetRiskLevelCollection(*mongoClient).FindOne(service.C, bson.D{{"_id", result.InsertedID}}).Decode(&createdRiskLevel)
	if err != nil {
		return riskLevel, err
	}

	service.InvalidatePublicCache()

	return &createdRiskLevel, nil
}

func (service *RiskLevelService) PublicRiskLevels() (riskLevels []*models.RiskLevel, err error) {
	// try to load risklevels from cache
	riskLevels, err = publicRiskLevelsFromCache(service.C)
	if err == nil {
		return riskLevels, nil
	}

	// else, load them from the DB and put them in the cache
	mongoClient := db.GetDbClient()
	filter := bson.M{"archivedAt": nil}
	opts := options.Find().SetSort(bson.D{{"createdAt", -1}})
	cur, err := collections.GetRiskLevelCollection(*mongoClient).Find(service.C, filter, opts)

	if err != nil {
		return riskLevels, err
	}

	err = cur.All(service.C, &riskLevels)
	if err != nil {
		return riskLevels, err
	}

	storePublicRiskLevelsInCache(service.C, riskLevels)

	return riskLevels, nil
}

func (service *RiskLevelService) InvalidatePublicCache() {
	cacheClient := db.GetCacheClient()
	cacheClient.Del(service.C, cacheKey)
}

func publicRiskLevelsFromCache(c *gin.Context) (riskLevels []*models.RiskLevel, err error) {
	cacheClient := db.GetCacheClient()
	riskLevelsString, err := cacheClient.Get(c, cacheKey).Result()
	if err != nil {
		return riskLevels, err
	}

	err = json.Unmarshal([]byte(riskLevelsString), &riskLevels)
	if err != nil {
		return riskLevels, err
	}

	return riskLevels, nil
}

func storePublicRiskLevelsInCache(c *gin.Context, riskLevels []*models.RiskLevel) {
	cacheClient := db.GetCacheClient()
	riskLevelsString, err := json.Marshal(riskLevels)
	if err != nil {
		log.Printf("Error storing public riskLevels in cache: %s", err.Error())
		return
	}
	cacheClient.Set(c, cacheKey, riskLevelsString, 0)
}
