package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"backend_coffepaste/config"
	"backend_coffepaste/routes"
)

func main() {
	_ = godotenv.Load()

	db, err := config.ConnectDatabase()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://localhost:5174",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
		AllowCredentials: true,
	}))

	routes.SetupRoutes(router, db)

	println("================================")
	println("       CAFE POS BACKEND")
	println("================================")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	println("Server : 0.0.0.0:" + port)
	println("Health : /api/v1/health")
	println("================================")

	err = router.Run("0.0.0.0:" + port)
	if err != nil {
		panic(err)
	}
}
