package services

import (
	"almanac-api/collections"
	"almanac-api/db"

	"github.com/gin-gonic/gin"
	"gitlab.com/almanac-app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GovernorateService struct {
	C *gin.Context
}

type UpdateGovernorate struct {
	RiskLevel int `json:"riskLevel"`
}

func (service *GovernorateService) Update(id primitive.ObjectID, UpdateGovernorate *UpdateGovernorate) (governorate *models.Governorate, err error) {
	mongoClient := db.GetDbClient()

	filter := bson.M{"_id": id}

	update := bson.M{"$set": bson.M{
		"riskLevel": UpdateGovernorate.RiskLevel,
	}}

	err = collections.GetGovernorateCollection(*mongoClient).FindOneAndUpdate(service.C, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&governorate)
	if err != nil {
		return governorate, err
	}

	rlService := RiskLevelService{C: service.C}
	rlService.InvalidatePublicCache()

	return governorate, nil
}
