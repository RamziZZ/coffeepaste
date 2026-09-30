package orders

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
	orders, err := h.service.GetAll(
		c.Request.Context(),
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
		"data":    orders,
	})
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(
		c.Param("id"),
		10,
		64,
	)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "id order tidak valid",
		})
		return
	}

	order, err := h.service.GetByID(
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
		"data":    order,
	})
}

func (h *Handler) GetMyOrders(c *gin.Context) {
	customerID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "user tidak terautentikasi",
		})
		return
	}

	id, ok := customerID.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "user id tidak valid",
		})
		return
	}

	orders, err := h.service.GetByCustomerID(
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
		"data":    orders,
	})
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	customerID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "user tidak terautentikasi",
		})
		return
	}

	id, ok := customerID.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "user id tidak valid",
		})
		return
	}

	order, err := h.service.Create(
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
		"message": "order berhasil dibuat",
		"data":    order,
	})
}

func (h *Handler) AssignCashier(c *gin.Context) {
	id, err := strconv.ParseInt(
		c.Param("id"),
		10,
		64,
	)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "id order tidak valid",
		})
		return
	}

	var req AssignCashierRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	order, err := h.service.AssignCashier(
		c.Request.Context(),
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
		"message": "kasir berhasil ditetapkan",
		"data":    order,
	})
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(
		c.Param("id"),
		10,
		64,
	)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "id order tidak valid",
		})
		return
	}

	var req UpdateOrderStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	order, err := h.service.UpdateStatus(
		c.Request.Context(),
		id,
		req.Status,
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
		"message": "status order berhasil diubah",
		"data":    order,
	})
}
