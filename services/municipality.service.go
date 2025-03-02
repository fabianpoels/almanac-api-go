package services

import (
	"almanac-api/collections"
	"almanac-api/db"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/almanac-app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MunicipalityService struct {
	C *gin.Context
}

type UpdateMunicipality struct {
	RiskLevel int `json:"riskLevel"`
}

func (service *MunicipalityService) Create(createMunicipality *models.Municipality) (municipality *models.Municipality, err error) {
	mongoClient := db.GetDbClient()

	createMunicipality.CreatedAt = time.Now()
	createMunicipality.UpdatedAt = time.Now()

	result, err := collections.GetMunicipalityCollection(*mongoClient).InsertOne(service.C, createMunicipality)
	if err != nil {
		return municipality, err
	}

	rlService := RiskLevelService{C: service.C}
	rlService.InvalidatePublicCache()

	err = collections.GetMunicipalityCollection(*mongoClient).FindOne(service.C, bson.D{{"_id", result.InsertedID}}).Decode(&municipality)
	if err != nil {
		return municipality, err
	}

	return municipality, nil
}

func (service *MunicipalityService) Update(id primitive.ObjectID, updateMunicipality *UpdateMunicipality) (municipality *models.Municipality, err error) {
	mongoClient := db.GetDbClient()

	filter := bson.M{"_id": id}

	update := bson.M{"$set": bson.M{
		"riskLevel": updateMunicipality.RiskLevel,
	}}

	err = collections.GetMunicipalityCollection(*mongoClient).FindOneAndUpdate(service.C, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&municipality)
	if err != nil {
		return municipality, err
	}

	rlService := RiskLevelService{C: service.C}
	rlService.InvalidatePublicCache()

	return municipality, nil
}

func (service *MunicipalityService) Delete(id primitive.ObjectID) (delete bool, err error) {
	mongoClient := db.GetDbClient()

	filter := bson.M{"_id": id, "osmId": ""}

	_, err = collections.GetMunicipalityCollection(*mongoClient).DeleteOne(service.C, filter)
	if err != nil {
		return false, err
	}

	riskLevelService := RiskLevelService{C: service.C}
	riskLevelService.InvalidatePublicCache()

	return true, nil
}
