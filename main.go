package main

import (
	"First_Go_Gorm/DB"
	"First_Go_Gorm/handlers"
	"First_Go_Gorm/middlewares"
	"github.com/gin-gonic/gin"
)

const SecretKey = "abc"

func main() {
	DB := DB.Init()
	h := handlers.New(DB)
	r := gin.Default()
	r.Use(gin.Logger())
	apiRoutes := r.Group("/api", middlewares.AuthJWT())
	{
		apiRoutes.GET("/List", h.GetAllTask)
		//apiRoutes.POST("/Add", h.AddTask)
		apiRoutes.PUT("/Update/:id", h.UpdateTask)
		apiRoutes.DELETE("/Delete/:id", h.DeleteTask)
	}
	r.POST("/Add", h.AddTask)
	r.POST("/auth/register", h.Register)
	r.POST("/auth/login", h.Login)
	go r.Run(":8080")
	send := gin.Default()
	send.POST("/", h.Senduser)
	send.Run(":8081")

}
