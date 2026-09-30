package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"backend_coffepaste/config"
	"backend_coffepaste/routes"
)

func main() {

	// Membaca file .env
	err := godotenv.Load()
	if err != nil {
		panic("Gagal membaca file .env")
	}

	// Koneksi database
	db, err := config.ConnectDatabase()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Membuat Gin
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, db)

	//Informasi Server
	println("================================")
	println("       CAFE POS BACKEND")
	println("================================")
	println("Server : http://localhost:8080")
	println("Health : http://localhost:8080/api/v1/health")
	println("================================")

	// Jalankan server
	router.Run(":8080")
}
