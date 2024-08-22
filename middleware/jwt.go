package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"

	"almanac-api/db"
	"almanac-api/models"
	"almanac-api/utils"
)

func ValidateJwt() gin.HandlerFunc {
	// load non-request based stuff
	var mongoClient = db.GetDbClient()

	return func(c *gin.Context) {
		// validate the JWT and read the user _id
		tokenString := c.GetHeader("Authorization")[len("Bearer "):]
		token, err := utils.ParseJwt(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Token invalid"})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Token invalid"})
			return
		}
		email := claims["email"]

		// look up the user
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var user models.User
		err = models.GetUserCollection(*mongoClient).FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(&user)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "user not found"})
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
