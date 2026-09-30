package shifts

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.RouterGroup, handler *Handler) {
	shiftRoutes := router.Group("/shifts")
	{
		shiftRoutes.GET("", handler.GetAll)
		shiftRoutes.GET("/:id", handler.GetByID)
		shiftRoutes.GET("/my", handler.GetMyShifts)
		shiftRoutes.GET("/my/open", handler.GetMyOpenShift)
		shiftRoutes.POST("/open", handler.Open)
		shiftRoutes.POST("/:id/close", handler.Close)
	}
}