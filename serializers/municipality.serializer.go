package serializers

import (
	"gitlab.com/almanac-app/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MunicipalityResponse struct {
	Id   primitive.ObjectID     `json:"id"`
	Name models.TranslatedField `json:"name"`
}

type MunicipalitySerializer struct {
	models.Municipality
}

type MunicipalitiesSerializer struct {
	Municipalities []models.Municipality
}

func (ser *MunicipalitySerializer) Response() MunicipalityResponse {
	response := MunicipalityResponse{
		Id:   ser.Id,
		Name: ser.Name,
	}
	return response
}

func (ser *MunicipalitiesSerializer) Response() []MunicipalityResponse {
	response := []MunicipalityResponse{}
	for _, municipality := range ser.Municipalities {
		serializer := MunicipalitySerializer{municipality}
		response = append(response, serializer.Response())
	}
	return response
}
