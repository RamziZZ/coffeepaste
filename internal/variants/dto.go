package variants

type CreateVariantRequest struct {
	ProductID       int64   `json:"product_id" binding:"required"`
	Name            string  `json:"name" binding:"required"`
	PriceAdjustment float64 `json:"price_adjustment"`
	IsActive        *bool   `json:"is_active"`
}

type UpdateVariantRequest struct {
	ProductID       int64   `json:"product_id" binding:"required"`
	Name            string  `json:"name" binding:"required"`
	PriceAdjustment float64 `json:"price_adjustment"`
	IsActive        *bool   `json:"is_active"`
}
