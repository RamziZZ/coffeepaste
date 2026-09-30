package inventory

import (
	"context"
	"errors"
	"fmt"

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
) ([]InventoryMovement, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT
			id,
			product_id,
			type,
			quantity,
			stock_before,
			stock_after,
			reference_id,
			reference_type,
			note,
			created_at
		FROM inventory_movements
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movements := []InventoryMovement{}

	for rows.Next() {
		var movement InventoryMovement
		var referenceID string
		var referenceType string
		var note string

		err := rows.Scan(
			&movement.ID,
			&movement.ProductID,
			&movement.Type,
			&movement.Quantity,
			&movement.StockBefore,
			&movement.StockAfter,
			&referenceID,
			&referenceType,
			&note,
			&movement.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if referenceID != "" {
			movement.ReferenceID = &referenceID
		}

		if referenceType != "" {
			movement.ReferenceType = &referenceType
		}

		if note != "" {
			movement.Note = &note
		}

		movements = append(movements, movement)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movements, nil
}

func (r *Repository) FindByID(
	ctx context.Context,
	id int64,
) (*InventoryMovement, error) {
	var movement InventoryMovement
	var referenceID string
	var referenceType string
	var note string

	err := r.DB.QueryRow(ctx, `
		SELECT
			id,
			product_id,
			type,
			quantity,
			stock_before,
			stock_after,
			reference_id,
			reference_type,
			note,
			created_at
		FROM inventory_movements
		WHERE id = $1
	`, id).Scan(
		&movement.ID,
		&movement.ProductID,
		&movement.Type,
		&movement.Quantity,
		&movement.StockBefore,
		&movement.StockAfter,
		&referenceID,
		&referenceType,
		&note,
		&movement.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	if referenceID != "" {
		movement.ReferenceID = &referenceID
	}

	if referenceType != "" {
		movement.ReferenceType = &referenceType
	}

	if note != "" {
		movement.Note = &note
	}

	return &movement, nil
}

func (r *Repository) Create(
	ctx context.Context,
	req CreateInventoryMovementRequest,
) (*InventoryMovement, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	var stockBefore int

	err = tx.QueryRow(ctx, `
		SELECT stock
		FROM products
		WHERE id = $1
		FOR UPDATE
	`, req.ProductID).Scan(&stockBefore)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("produk tidak ditemukan")
		}

		return nil, err
	}

	stockAfter := stockBefore

	switch req.Type {
	case "IN":
		stockAfter = stockBefore + req.Quantity
	case "OUT":
		stockAfter = stockBefore - req.Quantity
	case "ADJUSTMENT":
		stockAfter = req.Quantity
	default:
		return nil, fmt.Errorf("type inventory tidak valid: %s", req.Type)
	}

	if stockAfter < 0 {
		return nil, errors.New("stok tidak mencukupi")
	}

	_, err = tx.Exec(ctx, `
		UPDATE products
		SET
			stock = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`,
		stockAfter,
		req.ProductID,
	)

	if err != nil {
		return nil, err
	}

	var movement InventoryMovement
	var referenceID string
	var referenceType string
	var note string

	err = tx.QueryRow(ctx, `
		INSERT INTO inventory_movements (
			product_id,
			type,
			quantity,
			stock_before,
			stock_after,
			reference_id,
			reference_type,
			note
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING
			id,
			product_id,
			type,
			quantity,
			stock_before,
			stock_after,
			reference_id,
			reference_type,
			note,
			created_at
	`,
		req.ProductID,
		req.Type,
		req.Quantity,
		stockBefore,
		stockAfter,
		req.ReferenceID,
		req.ReferenceType,
		req.Note,
	).Scan(
		&movement.ID,
		&movement.ProductID,
		&movement.Type,
		&movement.Quantity,
		&movement.StockBefore,
		&movement.StockAfter,
		&referenceID,
		&referenceType,
		&note,
		&movement.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	if referenceID != "" {
		movement.ReferenceID = &referenceID
	}

	if referenceType != "" {
		movement.ReferenceType = &referenceType
	}

	if note != "" {
		movement.Note = &note
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &movement, nil
}
