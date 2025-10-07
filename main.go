package main

import (
	"bitresume/config"
	"bitresume/jobs"
	"bitresume/routes"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	// Try loading .env (optional)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables from Railway...")
	}

	// Initialize configs
	config.InitOAuth()
	config.InitDB()

	// Setup Gin
	r := gin.Default()
	r.Static("/uploads", "./uploads")

	// CORS setup
	corsConfig := cors.Config{
		AllowOrigins:     []string{"https://student-smart-hub.web.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	r.Use(cors.New(corsConfig))

	// Root route
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Server is running...")
	})

	// Register routes
	routes.RegisterRoutes(r)

	// Cron job (runs daily at 11:55 PM)
	c := cron.New(cron.WithSeconds())
	_, errCron := c.AddFunc("0 55 23 * * *", jobs.CallDailyTasksForAllDates)
	if errCron != nil {
		panic("Failed to schedule cron job: " + errCron.Error())
	}
	c.Start()

	// Get PORT (Railway provides this automatically)
	port := "6001"

	log.Println("Server starting on port:", port)
	r.Run(":" + port)
}
