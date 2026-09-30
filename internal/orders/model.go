package orders

import "time"

type Order struct {
	ID          int64      `json:"id"`
	InvoiceNo   string     `json:"invoice_no"`
	CustomerID  int64      `json:"customer_id"`
	CashierID   *int64     `json:"cashier_id"`
	ShiftID     *int64     `json:"shift_id"`
	OrderType   string     `json:"order_type"`
	TableNumber *string    `json:"table_number"`
	Subtotal    float64    `json:"subtotal"`
	Discount    float64    `json:"discount"`
	Tax         float64    `json:"tax"`
	TotalAmount float64    `json:"total_amount"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}