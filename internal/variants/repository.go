package variants

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


// GET ALL VARIANTS
func (r *Repository) FindAll(ctx context.Context) ([]ProductVariant, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT
			id,
			product_id,
			name,
			price_adjustment,
			is_active
		FROM product_variants
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	variants := []ProductVariant{}

	for rows.Next() {
		var variant ProductVariant

		err := rows.Scan(
			&variant.ID,
			&variant.ProductID,
			&variant.Name,
			&variant.PriceAdjustment,
			&variant.IsActive,
		)
		if err != nil {
			return nil, err
		}

		variants = append(variants, variant)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return variants, nil
}

// GET VARIANT BY ID
func (r *Repository) FindByID(
	ctx context.Context,
	id int64,
) (*ProductVariant, error) {
	variant := &ProductVariant{}

	err := r.DB.QueryRow(ctx, `
		SELECT
			id,
			product_id,
			name,
			price_adjustment,
			is_active
		FROM product_variants
		WHERE id = $1
	`, id).Scan(
		&variant.ID,
		&variant.ProductID,
		&variant.Name,
		&variant.PriceAdjustment,
		&variant.IsActive,
	)

	if err != nil {
		return nil, err
	}

	return variant, nil
}


// CREATE VARIANT
func (r *Repository) Create(
	ctx context.Context,
	req CreateVariantRequest,
) (*ProductVariant, error) {
	variant := &ProductVariant{}

	isActive := true

	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	err := r.DB.QueryRow(ctx, `
		INSERT INTO product_variants (
			product_id,
			name,
			price_adjustment,
			is_active
		)
		VALUES ($1, $2, $3, $4)
		RETURNING
			id,
			product_id,
			name,
			price_adjustment,
			is_active
	`,
		req.ProductID,
		req.Name,
		req.PriceAdjustment,
		isActive,
	).Scan(
		&variant.ID,
		&variant.ProductID,
		&variant.Name,
		&variant.PriceAdjustment,
		&variant.IsActive,
	)

	if err != nil {
		return nil, err
	}

	return variant, nil
}

// UPDATE VARIANT
func (r *Repository) Update(
	ctx context.Context,
	id int64,
	req UpdateVariantRequest,
) (*ProductVariant, error) {
	variant := &ProductVariant{}

	isActive := true

	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	err := r.DB.QueryRow(ctx, `
		UPDATE product_variants
		SET
			product_id = $1,
			name = $2,
			price_adjustment = $3,
			is_active = $4
		WHERE id = $5
		RETURNING
			id,
			product_id,
			name,
			price_adjustment,
			is_active
	`,
		req.ProductID,
		req.Name,
		req.PriceAdjustment,
		isActive,
		id,
	).Scan(
		&variant.ID,
		&variant.ProductID,
		&variant.Name,
		&variant.PriceAdjustment,
		&variant.IsActive,
	)

	if err != nil {
		return nil, err
	}

	return variant, nil
}


// DELETE VARIANT
func (r *Repository) Delete(
	ctx context.Context,
	id int64,
) error {
	_, err := r.DB.Exec(ctx, `
		DELETE FROM product_variants
		WHERE id = $1
	`, id)

	return err
}