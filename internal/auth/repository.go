package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r *Repository) FindByEmail(
	ctx context.Context,
	email string,
) (*User, error) {

	user := &User{}

	err := r.DB.QueryRow(
		ctx,
		`
		SELECT
			u.id,
			u.role_id,
			u.name,
			u.email,
			u.password_hash,
			u.is_active,
			r.name
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.email = $1
		`,
		email,
	).Scan(
		&user.ID,
		&user.RoleID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.IsActive,
		&user.Role,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) Create(
	ctx context.Context,
	name string,
	email string,
	passwordHash string,
) (*User, error) {

	user := &User{}

	err := r.DB.QueryRow(
		ctx,
		`
		INSERT INTO users
			(role_id, name, email, password_hash, is_active)
		VALUES
			(3, $1, $2, $3, true)
		RETURNING
			id,
			role_id,
			name,
			email,
			password_hash,
			is_active
		`,
		name,
		email,
		passwordHash,
	).Scan(
		&user.ID,
		&user.RoleID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.IsActive,
	)

	if err != nil {
		return nil, err
	}

	user.Role = "CUSTOMER"

	return user, nil
}
