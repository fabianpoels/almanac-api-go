package services

import (
	"almanac-api/collections"
	"almanac-api/db"
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"gitlab.com/almanac-app/models"
	"go.mongodb.org/mongo-driver/bson"
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

// SeverityLevel represents the different severity categories
type SeverityLevel string

// Define the possible severity levels
const (
	Minor    SeverityLevel = "minor"
	Moderate SeverityLevel = "moderate"
	Severe   SeverityLevel = "severe"
)

type RiskLevelResponse struct {
	Minor    models.GeoJSON `json:"minor"`
	Moderate models.GeoJSON `json:"moderate"`
	Severe   models.GeoJSON `json:"severe"`
}

func (service *RiskLevelService) PublicRiskLevels() (response RiskLevelResponse, err error) {
	// try to load risklevels from cache
	response, err = publicRiskLevelsFromCache(service.C)
	if err == nil {
		return response, nil
	}

	mongoClient := db.GetDbClient()

	// find municipalities with risk level >= 0
	filter := bson.M{"riskLevel": bson.M{"$gte": 0}}
	cur, err := collections.GetMunicipalityCollection(*mongoClient).Find(service.C, filter, options.Find())
	if err != nil {
		return response, err
	}

	municipalities := make([]models.Municipality, 0)
	err = cur.All(service.C, &municipalities)
	if err != nil {
		return response, err
	}

	response = RiskLevelResponse{
		Minor: models.GeoJSON{
			Type:     "FeatureCollection",
			Features: []models.GeoJSONFeature{},
		},
		Moderate: models.GeoJSON{
			Type:     "FeatureCollection",
			Features: []models.GeoJSONFeature{},
		},
		Severe: models.GeoJSON{
			Type:     "FeatureCollection",
			Features: []models.GeoJSONFeature{},
		},
	}
	for _, m := range municipalities {
		switch m.RiskLevel {
		case 0:
			response.Minor.Features = append(response.Minor.Features, m.GeoData.Features...)
		case 1:
			response.Moderate.Features = append(response.Moderate.Features, m.GeoData.Features...)
		case 2:
			response.Severe.Features = append(response.Severe.Features, m.GeoData.Features...)
		}
	}

	storePublicRiskLevelsInCache(service.C, response)

	return response, nil
}

func (service *RiskLevelService) InvalidatePublicCache() {
	cacheClient := db.GetCacheClient()
	cacheClient.Del(service.C, cacheKey)
}

func publicRiskLevelsFromCache(c *gin.Context) (response RiskLevelResponse, err error) {
	cacheClient := db.GetCacheClient()
	responseString, err := cacheClient.Get(c, cacheKey).Result()
	if err != nil {
		return response, err
	}

	err = json.Unmarshal([]byte(responseString), &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func storePublicRiskLevelsInCache(c *gin.Context, response RiskLevelResponse) {
	cacheClient := db.GetCacheClient()
	responseString, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error storing public riskLevels in cache: %s", err.Error())
		return
	}
	cacheClient.Set(c, cacheKey, responseString, 0)
}
