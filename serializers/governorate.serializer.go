package serializers

import (
	"gitlab.com/almanac-app/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GovernorateResponse struct {
	Id        primitive.ObjectID     `json:"id"`
	Name      models.TranslatedField `json:"name"`
	RiskLevel int                    `json:"riskLevel"`
	GeoData   models.GeoJSON         `json:"geoData" bson:"geoData"`
	OsmId     any                    `json:"osmId"`
}

type GovernorateSerializer struct {
	models.Governorate
}

type GovernoratesSerliazer struct {
	Governorates []models.Governorate
}

func (ser *GovernorateSerializer) Response() GovernorateResponse {
	response := GovernorateResponse{
		Id:        ser.Id,
		Name:      ser.Name,
		RiskLevel: ser.RiskLevel,
		GeoData:   ser.GeoData,
		OsmId:     ser.OsmID,
	}
	return response
}

func (ser *GovernoratesSerliazer) Response() []GovernorateResponse {
	response := []GovernorateResponse{}
	for _, governorate := range ser.Governorates {
		serializer := GovernorateSerializer{governorate}
		response = append(response, serializer.Response())
	}
	return response
}
