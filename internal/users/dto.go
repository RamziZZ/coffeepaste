package users

import "time"

type CreateCashierRequest struct {
	Name     string  `json:"name" binding:"required"`
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6"`
	Phone    *string `json:"phone"`
}

type UpdateUserRequest struct {
	Name  string  `json:"name" binding:"required"`
	Email string  `json:"email" binding:"required,email"`
	Phone *string `json:"phone"`
}

type UpdatePasswordRequest struct {
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateStatusRequest struct {
	IsActive bool `json:"is_active"`
}

type UserResponse struct {
	ID        int64     `json:"id"`
	RoleID    int       `json:"role_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     *string   `json:"phone"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}