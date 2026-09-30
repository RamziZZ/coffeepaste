package inventory

import "time"

type InventoryMovement struct {
	ID            int64     `json:"id"`
	ProductID     int64     `json:"product_id"`
	Type          string    `json:"type"`
	Quantity      int       `json:"quantity"`
	StockBefore   int       `json:"stock_before"`
	StockAfter    int       `json:"stock_after"`
	ReferenceID   *string   `json:"reference_id"`
	ReferenceType *string   `json:"reference_type"`
	Note          *string   `json:"note"`
	CreatedAt     time.Time `json:"created_at"`
}
