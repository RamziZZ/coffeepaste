package main

import (
	"backend_coffepaste/config"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http"
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

	// HEALTH CHECK
	router.GET("/api/v1/health", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Cafe POS Backend is running",
		})

	})

	// Informasi server
	println("================================")
	println("       CAFE POS BACKEND")
	println("================================")
	println("Server : http://localhost:8080")
	println("Health : http://localhost:8080/api/v1/health")
	println("================================")

	// Menjalankan server
	router.Run(":8080")
}
