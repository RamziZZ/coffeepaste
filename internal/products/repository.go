package products

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
) ([]Product, error) {
	rows, err := r.DB.Query(
		ctx,
		`
		SELECT
			id,
			category_id,
			name,
			COALESCE(barcode, ''),
			COALESCE(description, ''),
			price::bigint,
			stock,
			COALESCE(image_url, ''),
			is_available
		FROM products
		ORDER BY id ASC
		`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := []Product{}

	for rows.Next() {
		var product Product

		err := rows.Scan(
			&product.ID,
			&product.CategoryID,
			&product.Name,
			&product.Barcode,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.ImageURL,
			&product.IsAvailable,
		)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// GET PRODUCT BY ID
func (r *Repository) FindByID(
	ctx context.Context,
	id int64,
) (*Product, error) {
	product := &Product{}

	err := r.DB.QueryRow(
		ctx,
		`
		SELECT
			id,
			category_id,
			name,
			COALESCE(barcode, ''),
			COALESCE(description, ''),
			price::bigint,
			stock,
			COALESCE(image_url, ''),
			is_available
		FROM products
		WHERE id = $1
		`,
		id,
	).Scan(
		&product.ID,
		&product.CategoryID,
		&product.Name,
		&product.Barcode,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.ImageURL,
		&product.IsAvailable,
	)

	if err != nil {
		return nil, err
	}

	return product, nil
}

// CREATE PRODUCT
func (r *Repository) Create(
	ctx context.Context,
	req CreateProductRequest,
) (*Product, error) {
	product := &Product{}

	isAvailable := true

	if req.IsAvailable != nil {
		isAvailable = *req.IsAvailable
	}

	err := r.DB.QueryRow(
		ctx,
		`
		INSERT INTO products (
			category_id,
			name,
			barcode,
			description,
			price,
			stock,
			image_url,
			is_available
		)
		VALUES (
			$1,
			$2,
			NULLIF($3, ''),
			NULLIF($4, ''),
			$5,
			$6,
			NULLIF($7, ''),
			$8
		)
		RETURNING
			id,
			category_id,
			name,
			COALESCE(barcode, ''),
			COALESCE(description, ''),
			price::bigint,
			stock,
			COALESCE(image_url, ''),
			is_available
		`,
		req.CategoryID,
		req.Name,
		req.Barcode,
		req.Description,
		req.Price,
		req.Stock,
		req.ImageURL,
		isAvailable,
	).Scan(
		&product.ID,
		&product.CategoryID,
		&product.Name,
		&product.Barcode,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.ImageURL,
		&product.IsAvailable,
	)

	if err != nil {
		return nil, err
	}

	return product, nil
}

// UPDATE PRODUCT
func (r *Repository) Update(
	ctx context.Context,
	id int64,
	req UpdateProductRequest,
) (*Product, error) {
	product := &Product{}

	isAvailable := true

	if req.IsAvailable != nil {
		isAvailable = *req.IsAvailable
	}

	err := r.DB.QueryRow(
		ctx,
		`
		UPDATE products
		SET
			category_id = $1,
			name = $2,
			barcode = NULLIF($3, ''),
			description = NULLIF($4, ''),
			price = $5,
			stock = $6,
			image_url = NULLIF($7, ''),
			is_available = $8,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $9
		RETURNING
			id,
			category_id,
			name,
			COALESCE(barcode, ''),
			COALESCE(description, ''),
			price::bigint,
			stock,
			COALESCE(image_url, ''),
			is_available
		`,
		req.CategoryID,
		req.Name,
		req.Barcode,
		req.Description,
		req.Price,
		req.Stock,
		req.ImageURL,
		isAvailable,
		id,
	).Scan(
		&product.ID,
		&product.CategoryID,
		&product.Name,
		&product.Barcode,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.ImageURL,
		&product.IsAvailable,
	)

	if err != nil {
		return nil, err
	}

	return product, nil
}

// DELETE PRODUCT
func (r *Repository) Delete(
	ctx context.Context,
	id int64,
) error {
	_, err := r.DB.Exec(
		ctx,
		`
		DELETE FROM products
		WHERE id = $1
		`,
		id,
	)

	return err
}
