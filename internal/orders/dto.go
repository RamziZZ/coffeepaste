package orders

import "time"

type CreateOrderItemRequest struct {
	ProductID int64   `json:"product_id" binding:"required"`
	VariantID *int64  `json:"variant_id"`
	Quantity  int     `json:"quantity" binding:"required,min=1"`
	Notes     *string `json:"notes"`
}

type CreateOrderRequest struct {
	OrderType   string                   `json:"order_type" binding:"required"`
	TableNumber *string                  `json:"table_number"`
	Discount    float64                  `json:"discount"`
	Tax         float64                  `json:"tax"`
	Items       []CreateOrderItemRequest `json:"items" binding:"required,min=1"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type AssignCashierRequest struct {
	CashierID int64 `json:"cashier_id" binding:"required"`
	ShiftID   int64 `json:"shift_id" binding:"required"`
}

type OrderResponse struct {
	ID          int64     `json:"id"`
	InvoiceNo   string    `json:"invoice_no"`
	CustomerID  int64     `json:"customer_id"`
	CashierID   *int64    `json:"cashier_id"`
	ShiftID     *int64    `json:"shift_id"`
	OrderType   string    `json:"order_type"`
	TableNumber *string   `json:"table_number"`
	Subtotal    float64   `json:"subtotal"`
	Discount    float64   `json:"discount"`
	Tax         float64   `json:"tax"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type OrderItemResponse struct {
	ID          int64   `json:"id"`
	OrderID     int64   `json:"order_id"`
	ProductID   int64   `json:"product_id"`
	VariantID   *int64  `json:"variant_id"`
	ProductName string  `json:"product_name"`
	VariantName *string `json:"variant_name"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Subtotal    float64 `json:"subtotal"`
	Notes       *string `json:"notes"`
}

type OrderDetailResponse struct {
	Order OrderResponse       `json:"order"`
	Items []OrderItemResponse `json:"items"`
}
