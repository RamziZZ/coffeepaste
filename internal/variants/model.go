package variants

type ProductVariant struct {
	ID              int64   `json:"id"`
	ProductID       int64   `json:"product_id"`
	Name            string  `json:"name"`
	PriceAdjustment float64 `json:"price_adjustment"`
	IsActive        bool    `json:"is_active"`
}
