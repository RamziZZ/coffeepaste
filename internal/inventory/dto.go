package inventory

type CreateInventoryMovementRequest struct {
	ProductID     int64   `json:"product_id" binding:"required"`
	Type          string  `json:"type" binding:"required"`
	Quantity      int     `json:"quantity" binding:"required"`
	ReferenceID   *string `json:"reference_id"`
	ReferenceType *string `json:"reference_type"`
	Note          *string `json:"note"`
}
