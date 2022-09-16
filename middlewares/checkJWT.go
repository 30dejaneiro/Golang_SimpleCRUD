package middlewares

import (
	"First_Go_Gorm/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		//const BEARER_SCHERMA = ""
		//authHeader := c.GetHeader("Authorization")
		//tokenString := authHeader[len(BEARER_SCHERMA):]
		tokenString := c.GetHeader("Authorization")
		token, err := models.NewJWTService().ValidateToken(tokenString)
		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			log.Println("Claims[username]: ", claims["username"])
			log.Println("Claims[IsAdmin]: ", claims["isAdmin"])
			log.Println("Claims[IssuedAt]: ", claims["iat"])
			log.Println("Claims[ExpiresAt]: ", claims["exp"])
			if claims["isAdmin"] == true {
				c.Next()
			} else {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"err": "Cannot access",
				})
			}
		} else {
			log.Println(err)
			c.AbortWithStatusJSON(401, gin.H{
				"err": "Invalid token",
			})
		}
	}
}
