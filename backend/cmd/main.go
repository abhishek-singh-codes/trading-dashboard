package main

import (
	"log"

	"trading-dashboard/internal/auth"
	"trading-dashboard/internal/config"
	"trading-dashboard/internal/db"
	"trading-dashboard/internal/market"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	database := db.Connect(cfg)
	defer database.Close()
	db.Migrate(database)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	authHandler := auth.NewHandler(database, cfg)
	api := r.Group("/api")

	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)

	protected := api.Group("/")
	protected.Use(auth.JWTMiddleware(cfg))
	{
		protected.GET("/auth/me", authHandler.Me)
		mh := market.NewHandler(database)
		protected.GET("/market/search", mh.Search)
		protected.GET("/market/quote/:symbol", mh.GetQuote)
		protected.GET("/market/history/:symbol", mh.GetHistory)
	}

	log.Printf("🚀 Server running on http://localhost:%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
