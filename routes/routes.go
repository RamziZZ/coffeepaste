package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"backend_coffepaste/internal/auth"
	"backend_coffepaste/internal/categories"
	"backend_coffepaste/internal/inventory"
	"backend_coffepaste/internal/orders"
	"backend_coffepaste/internal/products"
	"backend_coffepaste/internal/shifts"
	"backend_coffepaste/internal/users"
	"backend_coffepaste/internal/variants"
)

func SetupRoutes(
	router *gin.Engine,
	db *pgxpool.Pool,
) {
	api := router.Group("/api/v1")

	// HEALTH
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Cafe POS Backend is running",
		})
	})

	// AUTH
	authRepository := auth.NewRepository(db)
	authService := auth.NewService(authRepository)
	authHandler := auth.NewHandler(authService)

	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}

	// CATEGORIES
	categoryRepository := categories.NewRepository(db)
	categoryService := categories.NewService(categoryRepository)
	categoryHandler := categories.NewHandler(categoryService)

	categoryRoutes := api.Group("/categories")
	{
		categoryRoutes.GET("/", categoryHandler.GetAll)
		categoryRoutes.GET("/:id", categoryHandler.GetById)
		categoryRoutes.POST("/", categoryHandler.Create)
		categoryRoutes.PUT("/:id", categoryHandler.Update)
		categoryRoutes.DELETE("/:id", categoryHandler.Delete)
	}

	// PRODUCTS
	productRepository := products.NewRepository(db)
	productService := products.NewService(productRepository)
	productHandler := products.NewHandler(productService)

	productRoutes := api.Group("/products")
	{
		productRoutes.GET("/", productHandler.GetAll)
		productRoutes.GET("/:id", productHandler.GetByID)
		productRoutes.POST("/", productHandler.Create)
		productRoutes.PUT("/:id", productHandler.Update)
		productRoutes.DELETE("/:id", productHandler.Delete)
	}

	// VARIANTS
	variantRepository := variants.NewRepository(db)
	variantService := variants.NewService(variantRepository)
	variantHandler := variants.NewHandler(variantService)

	variantRoutes := api.Group("/variants")
	{
		variantRoutes.GET("", variantHandler.GetAll)
		variantRoutes.GET("/:id", variantHandler.GetByID)
		variantRoutes.POST("", variantHandler.Create)
		variantRoutes.PUT("/:id", variantHandler.Update)
		variantRoutes.DELETE("/:id", variantHandler.Delete)
	}

	// INVENTORY
	inventoryRepository := inventory.NewRepository(db)
	inventoryService := inventory.NewService(inventoryRepository)
	inventoryHandler := inventory.NewHandler(inventoryService)

	inventoryRoutes := api.Group("/inventory")
	{
		inventoryRoutes.GET("", inventoryHandler.GetAll)
		inventoryRoutes.GET("/:id", inventoryHandler.GetByID)
		inventoryRoutes.POST("", inventoryHandler.Create)
	}

	// USERS
	userRepository := users.NewRepository(db)
	userService := users.NewService(userRepository)
	userHandler := users.NewHandler(userService)

	userRoutes := api.Group("/users")
	{
		userRoutes.GET("", userHandler.GetAll)
		userRoutes.GET("/:id", userHandler.GetByID)
		userRoutes.POST("/cashier", userHandler.CreateCashier)
		userRoutes.PUT("/:id", userHandler.Update)
		userRoutes.PUT("/:id/password", userHandler.UpdatePassword)
		userRoutes.PATCH("/:id/status", userHandler.UpdateStatus)
		userRoutes.DELETE("/:id", userHandler.Delete)
	}

	// SHIFTS
	shiftRepository := shifts.NewRepository(db)
	shiftService := shifts.NewService(shiftRepository)
	shiftHandler := shifts.NewHandler(shiftService)

	shiftRoutes := api.Group("/shifts")
	{
		shiftRoutes.GET("", shiftHandler.GetAll)
		shiftRoutes.GET("/:id", shiftHandler.GetByID)
		shiftRoutes.GET("/my", shiftHandler.GetMyShifts)
		shiftRoutes.GET("/my/open", shiftHandler.GetMyOpenShift)
		shiftRoutes.POST("/open", shiftHandler.Open)
		shiftRoutes.POST("/:id/close", shiftHandler.Close)
	}

	// ORDERS
	orderRepository := orders.NewRepository(db)
	orderService := orders.NewService(orderRepository)
	orderHandler := orders.NewHandler(orderService)

	orderRoutes := api.Group("/orders")
	{
		orderRoutes.GET("", orderHandler.GetAll)
		orderRoutes.GET("/my", orderHandler.GetMyOrders)
		orderRoutes.GET("/:id", orderHandler.GetByID)
		orderRoutes.POST("", orderHandler.Create)
		orderRoutes.PATCH("/:id/assign", orderHandler.AssignCashier)
		orderRoutes.PATCH("/:id/status", orderHandler.UpdateStatus)
	}

}
