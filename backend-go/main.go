package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	ConnectDatabase()

	r := gin.Default()
	r.Use(cors.Default())

	api := r.Group("/api")
	{
		api.POST("/shield/toggle", ToggleShield)

		goApi := api.Group("/go")
		{
			goApi.GET("/data", GetAllData)
			goApi.POST("/data", CreateData)
			goApi.PUT("/data/:id", UpdateData)
			goApi.POST("/data/seed", SeedData)
			goApi.DELETE("/data/clear", ClearData)
		}
	}

	log.Println("Go Backend starting on port 8081...")
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
