package shifts

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	FindAll(ctx context.Context) ([]Shift, error)
	FindByID(ctx context.Context, id int64) (*Shift, error)
	FindByUserID(ctx context.Context, userID int64) ([]Shift, error)
	FindOpenByUserID(ctx context.Context, userID int64) (*Shift, error)
	Open(ctx context.Context, userID int64, openingCash float64) (*Shift, error)
	Close(ctx context.Context, id int64, closingCash float64) (*Shift, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindAll(ctx context.Context) ([]Shift, error) {
	query := `
		SELECT
			id,
			user_id,
			opened_at,
			closed_at,
			opening_cash,
			closing_cash,
			status,
			created_at
		FROM shifts
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil semua shift: %w", err)
	}
	defer rows.Close()

	var shifts []Shift

	for rows.Next() {
		var shift Shift

		err := rows.Scan(
			&shift.ID,
			&shift.UserID,
			&shift.OpenedAt,
			&shift.ClosedAt,
			&shift.OpeningCash,
			&shift.ClosingCash,
			&shift.Status,
			&shift.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal membaca data shift: %w", err)
		}

		shifts = append(shifts, shift)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("gagal membaca hasil query shift: %w", err)
	}

	return shifts, nil
}

func (r *repository) FindByID(ctx context.Context, id int64) (*Shift, error) {
	query := `
		SELECT
			id,
			user_id,
			opened_at,
			closed_at,
			opening_cash,
			closing_cash,
			status,
			created_at
		FROM shifts
		WHERE id = $1
	`

	var shift Shift

	err := r.db.QueryRow(ctx, query, id).Scan(
		&shift.ID,
		&shift.UserID,
		&shift.OpenedAt,
		&shift.ClosedAt,
		&shift.OpeningCash,
		&shift.ClosingCash,
		&shift.Status,
		&shift.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("gagal mengambil shift: %w", err)
	}

	return &shift, nil
}

func (r *repository) FindByUserID(ctx context.Context, userID int64) ([]Shift, error) {
	query := `
		SELECT
			id,
			user_id,
			opened_at,
			closed_at,
			opening_cash,
			closing_cash,
			status,
			created_at
		FROM shifts
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil shift user: %w", err)
	}
	defer rows.Close()

	var shifts []Shift

	for rows.Next() {
		var shift Shift

		err := rows.Scan(
			&shift.ID,
			&shift.UserID,
			&shift.OpenedAt,
			&shift.ClosedAt,
			&shift.OpeningCash,
			&shift.ClosingCash,
			&shift.Status,
			&shift.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal membaca shift user: %w", err)
		}

		shifts = append(shifts, shift)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("gagal membaca hasil query shift user: %w", err)
	}

	return shifts, nil
}

func (r *repository) FindOpenByUserID(ctx context.Context, userID int64) (*Shift, error) {
	query := `
		SELECT
			id,
			user_id,
			opened_at,
			closed_at,
			opening_cash,
			closing_cash,
			status,
			created_at
		FROM shifts
		WHERE user_id = $1
		AND status = 'OPEN'
		ORDER BY created_at DESC
		LIMIT 1
	`

	var shift Shift

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&shift.ID,
		&shift.UserID,
		&shift.OpenedAt,
		&shift.ClosedAt,
		&shift.OpeningCash,
		&shift.ClosingCash,
		&shift.Status,
		&shift.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("gagal mengecek shift aktif: %w", err)
	}

	return &shift, nil
}

func (r *repository) Open(
	ctx context.Context,
	userID int64,
	openingCash float64,
) (*Shift, error) {
	query := `
		INSERT INTO shifts (
			user_id,
			opened_at,
			opening_cash,
			status
		)
		VALUES (
			$1,
			NOW(),
			$2,
			'OPEN'
		)
		RETURNING
			id,
			user_id,
			opened_at,
			closed_at,
			opening_cash,
			closing_cash,
			status,
			created_at
	`

	var shift Shift

	err := r.db.QueryRow(
		ctx,
		query,
		userID,
		openingCash,
	).Scan(
		&shift.ID,
		&shift.UserID,
		&shift.OpenedAt,
		&shift.ClosedAt,
		&shift.OpeningCash,
		&shift.ClosingCash,
		&shift.Status,
		&shift.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("gagal membuka shift: %w", err)
	}

	return &shift, nil
}

func (r *repository) Close(
	ctx context.Context,
	id int64,
	closingCash float64,
) (*Shift, error) {
	query := `
		UPDATE shifts
		SET
			closed_at = NOW(),
			closing_cash = $1,
			status = 'CLOSED'
		WHERE id = $2
		AND status = 'OPEN'
		RETURNING
			id,
			user_id,
			opened_at,
			closed_at,
			opening_cash,
			closing_cash,
			status,
			created_at
	`

	var shift Shift

	err := r.db.QueryRow(
		ctx,
		query,
		closingCash,
		id,
	).Scan(
		&shift.ID,
		&shift.UserID,
		&shift.OpenedAt,
		&shift.ClosedAt,
		&shift.OpeningCash,
		&shift.ClosingCash,
		&shift.Status,
		&shift.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("gagal menutup shift: %w", err)
	}

	return &shift, nil
}
