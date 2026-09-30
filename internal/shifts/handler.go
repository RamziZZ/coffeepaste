package shifts

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetAll(c *gin.Context) {
	shifts, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    shifts,
	})
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "id shift tidak valid",
		})
		return
	}

	shift, err := h.service.GetByID(
		c.Request.Context(),
		id,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    shift,
	})
}

func (h *Handler) GetMyShifts(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "user tidak terautentikasi",
		})
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "user id tidak valid",
		})
		return
	}

	shifts, err := h.service.GetByUserID(
		c.Request.Context(),
		int64(id),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    shifts,
	})
}

func (h *Handler) GetMyOpenShift(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "user tidak terautentikasi",
		})
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "user id tidak valid",
		})
		return
	}

	shift, err := h.service.GetOpenByUserID(
		c.Request.Context(),
		int64(id),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if shift == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    nil,
			"message": "tidak ada shift aktif",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    shift,
	})
}

func (h *Handler) Open(c *gin.Context) {
	var req OpenShiftRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "user tidak terautentikasi",
		})
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "user id tidak valid",
		})
		return
	}

	shift, err := h.service.Open(
		c.Request.Context(),
		int64(id),
		req,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "shift berhasil dibuka",
		"data":    shift,
	})
}

func (h *Handler) Close(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "id shift tidak valid",
		})
		return
	}

	var req CloseShiftRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "user tidak terautentikasi",
		})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "user id tidak valid",
		})
		return
	}

	shift, err := h.service.Close(
		c.Request.Context(),
		int64(userIDInt),
		id,
		req,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "shift berhasil ditutup",
		"data":    shift,
	})
}
