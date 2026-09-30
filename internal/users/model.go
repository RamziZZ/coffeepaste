package users

import "time"

type User struct {
	ID           int64     `json:"id"`
	RoleID       int       `json:"role_id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Phone        *string   `json:"phone"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
