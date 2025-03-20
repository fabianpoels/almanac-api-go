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

type MunicipalityRiskLevels struct {
	Minor    models.GeoJSON `json:"minor"`
	Moderate models.GeoJSON `json:"moderate"`
	Severe   models.GeoJSON `json:"severe"`
}

type RiskLevelResponse struct {
	Governorates   []models.Governorate   `json:"governorates"`
	Municipalities MunicipalityRiskLevels `json:"municipalities"`
}

func (service *RiskLevelService) PublicRiskLevels() (response RiskLevelResponse, err error) {
	// try to load risklevels from cache
	response, err = publicRiskLevelsFromCache(service.C)
	if err == nil {
		return response, nil
	}

	mongoClient := db.GetDbClient()

	// governorates with risklevel >= 0
	filter := bson.M{"riskLevel": bson.M{"$gte": 0}}
	cur, err := collections.GetGovernorateCollection(*mongoClient).Find(service.C, filter, options.Find())
	if err != nil {
		return response, err
	}

	governorates := make([]models.Governorate, 0)
	err = cur.All(service.C, &governorates)
	if err != nil {
		return response, err
	}

	governorateRiskLevels := make(map[string]int)
	for _, gov := range governorates {
		governorateRiskLevels[gov.OsmID] = gov.RiskLevel
	}

	// find municipalities with risk level >= 0
	cur, err = collections.GetMunicipalityCollection(*mongoClient).Find(service.C, filter, options.Find())
	if err != nil {
		return response, err
	}

	allMunicipalities := make([]models.Municipality, 0)
	err = cur.All(service.C, &allMunicipalities)
	if err != nil {
		return response, err
	}

	// filter municipalities to only include those with different risk levels from their governorates
	municipalities := make([]models.Municipality, 0)
	for _, m := range allMunicipalities {
		if govRiskLevel, exists := governorateRiskLevels[m.GovernorateOsmID]; exists {
			if m.RiskLevel != govRiskLevel {
				municipalities = append(municipalities, m)
			}
		} else {
			municipalities = append(municipalities, m)
		}
	}

	response = RiskLevelResponse{
		Governorates: governorates,
		Municipalities: MunicipalityRiskLevels{
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
		},
	}

	for _, m := range municipalities {
		switch m.RiskLevel {
		case 0:
			response.Municipalities.Minor.Features = append(response.Municipalities.Minor.Features, m.GeoData.Features...)
		case 1:
			response.Municipalities.Moderate.Features = append(response.Municipalities.Moderate.Features, m.GeoData.Features...)
		case 2:
			response.Municipalities.Severe.Features = append(response.Municipalities.Severe.Features, m.GeoData.Features...)
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
