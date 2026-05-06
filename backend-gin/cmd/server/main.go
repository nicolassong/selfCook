package main

import (
	"log"
	"os"

	gin "github.com/gin-gonic/gin"

	"github.com/gin-gonic/gin/internal/app"
)

func main() {
	cfg := app.LoadConfig()

	database, err := app.OpenMySQL(cfg)
	if err != nil {
		log.Fatalf("open mysql failed: %v", err)
	}
	defer func() {
		sqlDB, err := database.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}()

	router := gin.Default()
	app.RegisterRoutes(router, database, cfg)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("start server failed: %v", err)
	}
}
