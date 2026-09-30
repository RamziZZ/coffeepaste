package shifts

import "time"

type Shift struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	OpenedAt    time.Time  `json:"opened_at"`
	ClosedAt    *time.Time `json:"closed_at"`
	OpeningCash float64    `json:"opening_cash"`
	ClosingCash *float64   `json:"closing_cash"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
}