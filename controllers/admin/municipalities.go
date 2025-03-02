package admin

import (
	"almanac-api/collections"
	"almanac-api/db"
	"almanac-api/middleware"
	"almanac-api/serializers"
	"almanac-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/almanac-app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MunicipalitiesController struct {
}

func (m MunicipalitiesController) List(c *gin.Context) {
	mongoClient := db.GetDbClient()
	cur, err := collections.GetMunicipalityCollection(*mongoClient).Find(c, bson.D{{}})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	municipalities := make([]models.Municipality, 0)
	err = cur.All(c, &municipalities)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	serializer := serializers.MunicipalitiesSerializer{
		Municipalities: municipalities,
	}
	c.JSON(http.StatusOK, serializer.Response())
}

func (m MunicipalitiesController) Create(c *gin.Context) {
	user, ok := middleware.GetUserFromContext(c)
	if !ok || !user.IsAdmin() {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Unauthorized": "Not admin"})
		return
	}

	var createMunicipality models.Municipality
	err := c.BindJSON(&createMunicipality)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid request body"})
		return
	}

	mService := services.MunicipalityService{C: c}
	municipality, err := mService.Create(&createMunicipality)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	serializer := serializers.MunicipalitySerializer{
		Municipality: *municipality,
	}

	c.JSON(http.StatusOK, serializer.Response())
}

func (m MunicipalitiesController) Update(c *gin.Context) {
	user, ok := middleware.GetUserFromContext(c)
	if !ok || !user.IsAdmin() {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Unauthorized": "Not admin"})
		return
	}

	id := c.Param("id")
	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "no valid id provided"})
		return
	}

	var updateMunicipality services.UpdateMunicipality
	err = c.BindJSON(&updateMunicipality)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid request body"})
		return
	}

	mService := services.MunicipalityService{C: c}
	municipality, err := mService.Update(objId, &updateMunicipality)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	serializer := serializers.MunicipalitySerializer{
		Municipality: *municipality,
	}

	c.JSON(http.StatusOK, serializer.Response())
}

func (m MunicipalitiesController) Delete(c *gin.Context) {
	id := c.Param("id")
	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "no valid id provided"})
		return
	}

	mService := services.MunicipalityService{C: c}
	_, err = mService.Delete(objId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bson.M{})
}
