package users

import (
	"context"

	"github.com/jackc/pgx/v5"
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

func (r *Repository) FindAll(
	ctx context.Context,
) ([]User, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT
			id,
			role_id,
			name,
			email,
			password_hash,
			phone,
			is_active,
			created_at,
			updated_at
		FROM users
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.ID,
			&user.RoleID,
			&user.Name,
			&user.Email,
			&user.PasswordHash,
			&user.Phone,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) FindByID(
	ctx context.Context,
	id int64,
) (*User, error) {
	var user User

	err := r.DB.QueryRow(ctx, `
		SELECT
			id,
			role_id,
			name,
			email,
			password_hash,
			phone,
			is_active,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`, id).Scan(
		&user.ID,
		&user.RoleID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Phone,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindByEmail(
	ctx context.Context,
	email string,
) (*User, error) {
	var user User

	err := r.DB.QueryRow(ctx, `
		SELECT
			id,
			role_id,
			name,
			email,
			password_hash,
			phone,
			is_active,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
	`, email).Scan(
		&user.ID,
		&user.RoleID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Phone,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) CreateCashier(
	ctx context.Context,
	name string,
	email string,
	passwordHash string,
	phone *string,
) (*User, error) {
	var user User

	err := r.DB.QueryRow(ctx, `
		INSERT INTO users (
			role_id,
			name,
			email,
			password_hash,
			phone,
			is_active
		)
		VALUES (
			2,
			$1,
			$2,
			$3,
			$4,
			true
		)
		RETURNING
			id,
			role_id,
			name,
			email,
			password_hash,
			phone,
			is_active,
			created_at,
			updated_at
	`,
		name,
		email,
		passwordHash,
		phone,
	).Scan(
		&user.ID,
		&user.RoleID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Phone,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) Update(
	ctx context.Context,
	id int64,
	req UpdateUserRequest,
) (*User, error) {
	var user User

	err := r.DB.QueryRow(ctx, `
		UPDATE users
		SET
			name = $1,
			email = $2,
			phone = $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING
			id,
			role_id,
			name,
			email,
			password_hash,
			phone,
			is_active,
			created_at,
			updated_at
	`,
		req.Name,
		req.Email,
		req.Phone,
		id,
	).Scan(
		&user.ID,
		&user.RoleID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Phone,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) UpdatePassword(
	ctx context.Context,
	id int64,
	passwordHash string,
) error {
	result, err := r.DB.Exec(ctx, `
		UPDATE users
		SET
			password_hash = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`,
		passwordHash,
		id,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *Repository) UpdateStatus(
	ctx context.Context,
	id int64,
	isActive bool,
) (*User, error) {
	var user User

	err := r.DB.QueryRow(ctx, `
		UPDATE users
		SET
			is_active = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING
			id,
			role_id,
			name,
			email,
			password_hash,
			phone,
			is_active,
			created_at,
			updated_at
	`,
		isActive,
		id,
	).Scan(
		&user.ID,
		&user.RoleID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Phone,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) Delete(
	ctx context.Context,
	id int64,
) error {
	if id == 1 {
		return pgx.ErrNoRows
	}

	result, err := r.DB.Exec(ctx, `
		DELETE FROM users
		WHERE id = $1
	`,
		id,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}