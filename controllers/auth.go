package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"almanac-api/config"
	"almanac-api/db"
	"almanac-api/middleware"
	"almanac-api/models"
	"almanac-api/utils"
)

const cookieName = "refreshToken"

type AuthController struct {
}

type UserLogin struct {
	Email    string `json:"email" bson:"email" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}

func (a AuthController) Login(c *gin.Context) {
	mongoClient := db.GetDbClient()
	cacheClient := db.GetCacheClient()
	var userLogin UserLogin

	if err := c.ShouldBindJSON(&userLogin); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := models.GetUserCollection(*mongoClient).FindOne(c, bson.D{{Key: "email", Value: userLogin.Email}}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "error")
		return
	}

	match, err := utils.VerifyPasswordHash(userLogin.Password, user.Password)

	if !match || err != nil {
		c.JSON(http.StatusUnauthorized, "error")
		return
	}

	jwtToken, err := utils.GenerateJwt(user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	maxAge, err := strconv.Atoi(config.GetConfig().GetString("refreshToken.maxAge"))
	if err != nil {
		maxAge = 86400
	}

	// generate refresh token
	refresh := utils.GenerateRefreshTokenString()

	// store refresh token and reverse in cache
	err1 := cacheClient.Do(c, cacheClient.B().Set().Key(refresh).Value(user.Id.String()).Nx().ExSeconds(int64(maxAge)).Build()).Error()
	err2 := cacheClient.Do(c, cacheClient.B().Set().Key(user.Id.String()).Value(refresh).Nx().ExSeconds(int64(maxAge)).Build()).Error()

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusUnauthorized, "error")
		return
	}

	// TODO: properly configure the cookie
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(cookieName, refresh, maxAge, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"user": user, "jwt": jwtToken})
}

func (a AuthController) RefreshToken(c *gin.Context) {
	cacheClient := db.GetCacheClient()
	mongoClient := db.GetDbClient()

	// get the cookie valie
	cookie, err := c.Cookie(cookieName)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// get the user id from the cache
	idString, err := cacheClient.Do(c, cacheClient.B().Get().Key(cookie).Build()).ToString()

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// convert the id from string to ObjectID
	objectId, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// look up the user in the db
	var user models.User
	err = models.GetUserCollection(*mongoClient).FindOne(c, bson.D{{Key: "_id", Value: objectId}}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "error")
		return
	}

	jwtToken, err := utils.GenerateJwt(user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user, "jwt": jwtToken})
}

func (a AuthController) Logout(c *gin.Context) {
	cacheClient := db.GetCacheClient()
	user, exists := middleware.GetUserFromContext(c)
	if !exists {
		log.Println("/logout without logged in user")
		c.JSON(http.StatusOK, "")
		return
	}

	refreshTokenString, err := cacheClient.Do(c, cacheClient.B().Get().Key(user.Id.String()).Build()).ToString()
	if err != nil {
		log.Printf("/logout user id not found in cache: %s", user.Id.String())
		c.JSON(http.StatusOK, "")
		return
	}

	err = cacheClient.Do(c, cacheClient.B().Del().Key(user.Id.String()).Build()).Error()
	if err != nil {
		log.Printf("/logout error deleting user id from cache: %s", user.Id.String())
		c.JSON(http.StatusOK, "")
		return
	}

	err = cacheClient.Do(c, cacheClient.B().Del().Key(refreshTokenString).Build()).Error()
	if err != nil {
		log.Printf("/logout error deleting refreshtoken from cache: %s", refreshTokenString)
		c.JSON(http.StatusOK, "")
		return
	}

	c.SetCookie(cookieName, "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{})
}

// func (a AuthController) Register(c *gin.Context) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	var user models.User
// 	defer cancel()

// 	if err := c.BindJSON(&user); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	password, err := utils.HashPassword(user.Password)

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 	}

// 	newUser := models.User{
// 		Id:       primitive.NewObjectID(),
// 		Email:    user.Email,
// 		Name:     user.Name,
// 		Password: password,
// 	}

// 	result, err := models.GetUserCollection(*client).InsertOne(ctx, newUser)

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 	}

// 	c.JSON(http.StatusCreated, result)
// }
