package main

import (
	"os"

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
