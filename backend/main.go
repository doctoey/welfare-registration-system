package main

import (
	"fmt"
	"log"
	"welfare-registration-backend/internal/controllers"
	"welfare-registration-backend/internal/repositories"
	"welfare-registration-backend/internal/usecases"

	"github.com/gin-gonic/gin"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(corsMiddleware())

	repo, err := repositories.NewSQLiteRepository("welfare.db")
	if err != nil {
		log.Fatalf("unable to connect Database: %v", err)
	}

	usecase := usecases.NewApplicationUsecase(repo)

	ctrl := controllers.NewApplicationController(usecase)

	ctrl.RegisterRoutes(r)

	port := ":8080"
	fmt.Printf("run on http://localhost%s\n", port)
	r.Run(port)
}
