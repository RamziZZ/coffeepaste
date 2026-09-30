package categories

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

// GET ALL
func (r *Repository) FindAll(
	ctx context.Context,
) ([]Category, error) {
	rows, err := r.DB.Query(
		ctx,
		`
		SELECT
			id,
			name,
			description
		FROM categories
		ORDER BY id ASC
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []Category{}

	for rows.Next() {
		var category Category

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

// GET BY ID
func (r *Repository) FindByID(
	ctx context.Context,
	id int,
) (*Category, error) {
	category := &Category{}

	err := r.DB.QueryRow(
		ctx,
		`
		SELECT
			id,
			name,
			description
		FROM categories
		WHERE id = $1
		`,
		id,
	).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
	)

	if err != nil {
		return nil, err
	}

	return category, nil
}

// CREATE
func (r *Repository) Create(
	ctx context.Context,
	req CreateCategoryRequest,
) (*Category, error) {
	category := &Category{}

	err := r.DB.QueryRow(
		ctx,
		`
		INSERT INTO categories
			(name, description)
		VALUES
			($1, $2)
		RETURNING
			id,
			name,
			description
		`,
		req.Name,
		req.Description,
	).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
	)

	if err != nil {
		return nil, err
	}

	return category, nil
}

// UPDATE
func (r *Repository) Update(
	ctx context.Context,
	id int,
	req UpdateCategoryRequest,
) (*Category, error) {
	category := &Category{}

	err := r.DB.QueryRow(
		ctx,
		`
		UPDATE categories
		SET
			name = $1,
			description = $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING
			id,
			name,
			description
		`,
		req.Name,
		req.Description,
		id,
	).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
	)

	if err != nil {
		return nil, err
	}

	return category, nil
}

// DELETE
func (r *Repository) Delete(
	ctx context.Context,
	id int,
) error {
	_, err := r.DB.Exec(
		ctx,
		`
		DELETE FROM categories
		WHERE id = $1
		`,
		id,
	)

	return err
}
