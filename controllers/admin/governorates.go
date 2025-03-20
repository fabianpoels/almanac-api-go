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

type GovernoratesController struct {
}

func (g GovernoratesController) List(c *gin.Context) {
	mongoClient := db.GetDbClient()
	cur, err := collections.GetGovernorateCollection(*mongoClient).Find(c, bson.D{{}})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	governorates := make([]models.Governorate, 0)
	err = cur.All(c, &governorates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	serializer := serializers.GovernoratesSerliazer{
		Governorates: governorates,
	}
	c.JSON(http.StatusOK, serializer.Response())
}

func (g GovernoratesController) Update(c *gin.Context) {
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

	var updateGovernorate services.UpdateGovernorate
	err = c.BindJSON(&updateGovernorate)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid request body"})
		return
	}

	gService := services.GovernorateService{C: c}
	governorate, err := gService.Update(objId, &updateGovernorate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	serializer := serializers.GovernorateSerializer{
		Governorate: *governorate,
	}

	c.JSON(http.StatusOK, serializer.Response())
}
