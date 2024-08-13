package main

import (
	"fmt"
	"mass/database"
	"mass/services"
	"mass/routes"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	go func() {
		port := os.Getenv("PORT1")
		fmt.Println("PORT ->", port)
		
		router := gin.New()
		router.Use(gin.Logger())
		router.Use(cors.Default())
		routes.Routes(router)
		router.Run(":" + port)
	}()

	go func() {
		// helpers.Cron()
		services.Cron()
	}()

	select {}
	// if port == "" {
	// 	port = "8000"routes
	// }

	// router := gin.New()
	// router.Use(gin.Logger())
	// routes.UserRoutes(router)
	// router.Use(middleware.Authentication())

	// routes.FoodRoutes(router)
	// routes.MenuRoutes(router)
	// routes.TableRoutes(router)
	// routes.OrderRoutes(router)
	// routes.OrderItemRoutes(router)
	// routes.InvoiceRoutes(router)

	// router.Run(":" + port)
}
