package auth

type User struct {
	ID           int
	RoleID       int
	Name         string
	Email        string
	PasswordHash string
	IsActive     bool
	Role         string
}
