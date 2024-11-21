package serializers

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/almanac-app/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MunicipalityResponse struct {
	Id   primitive.ObjectID     `json:"id"`
	Name models.TranslatedField `json:"name"`
}

type MunicipalitySerializer struct {
	C *gin.Context
	models.Municipality
}

type MunicipalitiesSerializer struct {
	C              *gin.Context
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
		serializer := MunicipalitySerializer{ser.C, municipality}
		response = append(response, serializer.Response())
	}
	return response
}
