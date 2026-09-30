package orders

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	FindAll(ctx context.Context) ([]Order, error)
	FindByID(ctx context.Context, id int64) (*OrderDetailResponse, error)
	FindByCustomerID(ctx context.Context, customerID int64) ([]Order, error)
	Create(ctx context.Context, customerID int64, invoiceNo string, req CreateOrderRequest) (*OrderDetailResponse, error)
	AssignCashier(ctx context.Context, id int64, cashierID int64, shiftID int64) (*Order, error)
	UpdateStatus(ctx context.Context, id int64, status string) (*Order, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindAll(ctx context.Context) ([]Order, error) {
	query := `
		SELECT
			id,
			invoice_no,
			customer_id,
			cashier_id,
			shift_id,
			order_type,
			table_number,
			subtotal,
			discount,
			tax,
			total_amount,
			status,
			created_at,
			updated_at
		FROM orders
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil semua order: %w", err)
	}
	defer rows.Close()

	var orders []Order

	for rows.Next() {
		var order Order

		err := rows.Scan(
			&order.ID,
			&order.InvoiceNo,
			&order.CustomerID,
			&order.CashierID,
			&order.ShiftID,
			&order.OrderType,
			&order.TableNumber,
			&order.Subtotal,
			&order.Discount,
			&order.Tax,
			&order.TotalAmount,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal membaca order: %w", err)
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("gagal membaca hasil order: %w", err)
	}

	if orders == nil {
		orders = []Order{}
	}

	return orders, nil
}

func (r *repository) FindByID(
	ctx context.Context,
	id int64,
) (*OrderDetailResponse, error) {
	orderQuery := `
		SELECT
			id,
			invoice_no,
			customer_id,
			cashier_id,
			shift_id,
			order_type,
			table_number,
			subtotal,
			discount,
			tax,
			total_amount,
			status,
			created_at,
			updated_at
		FROM orders
		WHERE id = $1
	`

	var order Order

	err := r.db.QueryRow(
		ctx,
		orderQuery,
		id,
	).Scan(
		&order.ID,
		&order.InvoiceNo,
		&order.CustomerID,
		&order.CashierID,
		&order.ShiftID,
		&order.OrderType,
		&order.TableNumber,
		&order.Subtotal,
		&order.Discount,
		&order.Tax,
		&order.TotalAmount,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("gagal mengambil order: %w", err)
	}

	items, err := r.findItemsByOrderID(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	return &OrderDetailResponse{
		Order: OrderResponse{
			ID:          order.ID,
			InvoiceNo:   order.InvoiceNo,
			CustomerID:  order.CustomerID,
			CashierID:   order.CashierID,
			ShiftID:     order.ShiftID,
			OrderType:   order.OrderType,
			TableNumber: order.TableNumber,
			Subtotal:    order.Subtotal,
			Discount:    order.Discount,
			Tax:         order.Tax,
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
			CreatedAt:   order.CreatedAt,
			UpdatedAt:   order.UpdatedAt,
		},
		Items: items,
	}, nil
}

func (r *repository) FindByCustomerID(
	ctx context.Context,
	customerID int64,
) ([]Order, error) {
	query := `
		SELECT
			id,
			invoice_no,
			customer_id,
			cashier_id,
			shift_id,
			order_type,
			table_number,
			subtotal,
			discount,
			tax,
			total_amount,
			status,
			created_at,
			updated_at
		FROM orders
		WHERE customer_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil order customer: %w", err)
	}
	defer rows.Close()

	var orders []Order

	for rows.Next() {
		var order Order

		err := rows.Scan(
			&order.ID,
			&order.InvoiceNo,
			&order.CustomerID,
			&order.CashierID,
			&order.ShiftID,
			&order.OrderType,
			&order.TableNumber,
			&order.Subtotal,
			&order.Discount,
			&order.Tax,
			&order.TotalAmount,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal membaca order customer: %w", err)
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("gagal membaca hasil order customer: %w", err)
	}

	if orders == nil {
		orders = []Order{}
	}

	return orders, nil
}

func (r *repository) Create(
	ctx context.Context,
	customerID int64,
	invoiceNo string,
	req CreateOrderRequest,
) (*OrderDetailResponse, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal memulai transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var subtotal float64

	type itemData struct {
		productID    int64
		variantID    *int64
		productName  string
		variantName  *string
		quantity     int
		unitPrice    float64
		itemSubtotal float64
		notes        *string
	}

	items := make([]itemData, 0, len(req.Items))

	for _, item := range req.Items {
		var productName string
		var productPrice float64
		var isAvailable bool

		err := tx.QueryRow(
			ctx,
			`
			SELECT
				name,
				price,
				is_available
			FROM products
			WHERE id = $1
			`,
			item.ProductID,
		).Scan(
			&productName,
			&productPrice,
			&isAvailable,
		)

		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, fmt.Errorf(
					"produk dengan id %d tidak ditemukan",
					item.ProductID,
				)
			}

			return nil, fmt.Errorf("gagal mengambil produk: %w", err)
		}

		if !isAvailable {
			return nil, fmt.Errorf(
				"produk dengan id %d tidak tersedia",
				item.ProductID,
			)
		}

		unitPrice := productPrice
		var variantName *string

		if item.VariantID != nil {
			var adjustment float64
			var name string
			var active bool

			err := tx.QueryRow(
				ctx,
				`
				SELECT
					name,
					price_adjustment,
					is_active
				FROM product_variants
				WHERE id = $1
				AND product_id = $2
				`,
				*item.VariantID,
				item.ProductID,
			).Scan(
				&name,
				&adjustment,
				&active,
			)

			if err != nil {
				if err == pgx.ErrNoRows {
					return nil, fmt.Errorf(
						"variant dengan id %d tidak ditemukan untuk produk %d",
						*item.VariantID,
						item.ProductID,
					)
				}

				return nil, fmt.Errorf("gagal mengambil variant: %w", err)
			}

			if !active {
				return nil, fmt.Errorf(
					"variant dengan id %d tidak aktif",
					*item.VariantID,
				)
			}

			unitPrice += adjustment
			variantName = &name
		}

		itemSubtotal := unitPrice * float64(item.Quantity)
		subtotal += itemSubtotal

		items = append(items, itemData{
			productID:    item.ProductID,
			variantID:    item.VariantID,
			productName:  productName,
			variantName:  variantName,
			quantity:     item.Quantity,
			unitPrice:    unitPrice,
			itemSubtotal: itemSubtotal,
			notes:        item.Notes,
		})
	}

	totalAmount := subtotal - req.Discount + req.Tax

	if totalAmount < 0 {
		return nil, fmt.Errorf("total order tidak boleh negatif")
	}

	var order Order

	err = tx.QueryRow(
		ctx,
		`
		INSERT INTO orders (
			invoice_no,
			customer_id,
			order_type,
			table_number,
			subtotal,
			discount,
			tax,
			total_amount,
			status
		)
		VALUES (
			$1,
			$2,
			$3::order_type,
			$4,
			$5,
			$6,
			$7,
			$8,
			'PENDING'::order_status
		)
		RETURNING
			id,
			invoice_no,
			customer_id,
			cashier_id,
			shift_id,
			order_type,
			table_number,
			subtotal,
			discount,
			tax,
			total_amount,
			status,
			created_at,
			updated_at
		`,
		invoiceNo,
		customerID,
		req.OrderType,
		req.TableNumber,
		subtotal,
		req.Discount,
		req.Tax,
		totalAmount,
	).Scan(
		&order.ID,
		&order.InvoiceNo,
		&order.CustomerID,
		&order.CashierID,
		&order.ShiftID,
		&order.OrderType,
		&order.TableNumber,
		&order.Subtotal,
		&order.Discount,
		&order.Tax,
		&order.TotalAmount,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("gagal membuat order: %w", err)
	}

	for _, item := range items {
		_, err = tx.Exec(
			ctx,
			`
			INSERT INTO order_items (
				order_id,
				product_id,
				variant_id,
				product_name,
				variant_name,
				quantity,
				unit_price,
				subtotal,
				notes
			)
			VALUES (
				$1,
				$2,
				$3,
				$4,
				$5,
				$6,
				$7,
				$8,
				$9
			)
			`,
			order.ID,
			item.productID,
			item.variantID,
			item.productName,
			item.variantName,
			item.quantity,
			item.unitPrice,
			item.itemSubtotal,
			item.notes,
		)

		if err != nil {
			return nil, fmt.Errorf("gagal membuat order item: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("gagal commit order: %w", err)
	}

	return r.FindByID(ctx, order.ID)
}

func (r *repository) AssignCashier(
	ctx context.Context,
	id int64,
	cashierID int64,
	shiftID int64,
) (*Order, error) {
	query := `
		UPDATE orders
		SET
			cashier_id = $1,
			shift_id = $2,
			updated_at = NOW()
		WHERE id = $3
		RETURNING
			id,
			invoice_no,
			customer_id,
			cashier_id,
			shift_id,
			order_type,
			table_number,
			subtotal,
			discount,
			tax,
			total_amount,
			status,
			created_at,
			updated_at
	`

	var order Order

	err := r.db.QueryRow(
		ctx,
		query,
		cashierID,
		shiftID,
		id,
	).Scan(
		&order.ID,
		&order.InvoiceNo,
		&order.CustomerID,
		&order.CashierID,
		&order.ShiftID,
		&order.OrderType,
		&order.TableNumber,
		&order.Subtotal,
		&order.Discount,
		&order.Tax,
		&order.TotalAmount,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("gagal menetapkan kasir: %w", err)
	}

	return &order, nil
}

func (r *repository) UpdateStatus(
	ctx context.Context,
	id int64,
	status string,
) (*Order, error) {
	query := `
		UPDATE orders
		SET
			status = $1::order_status,
			updated_at = NOW()
		WHERE id = $2
		RETURNING
			id,
			invoice_no,
			customer_id,
			cashier_id,
			shift_id,
			order_type,
			table_number,
			subtotal,
			discount,
			tax,
			total_amount,
			status,
			created_at,
			updated_at
	`

	var order Order

	err := r.db.QueryRow(
		ctx,
		query,
		status,
		id,
	).Scan(
		&order.ID,
		&order.InvoiceNo,
		&order.CustomerID,
		&order.CashierID,
		&order.ShiftID,
		&order.OrderType,
		&order.TableNumber,
		&order.Subtotal,
		&order.Discount,
		&order.Tax,
		&order.TotalAmount,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("gagal mengubah status order: %w", err)
	}

	return &order, nil
}

func (r *repository) findItemsByOrderID(
	ctx context.Context,
	orderID int64,
) ([]OrderItemResponse, error) {
	query := `
		SELECT
			id,
			order_id,
			product_id,
			variant_id,
			product_name,
			variant_name,
			quantity,
			unit_price,
			subtotal,
			notes
		FROM order_items
		WHERE order_id = $1
		ORDER BY id ASC
	`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil order items: %w", err)
	}
	defer rows.Close()

	var items []OrderItemResponse

	for rows.Next() {
		var item OrderItemResponse

		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.VariantID,
			&item.ProductName,
			&item.VariantName,
			&item.Quantity,
			&item.UnitPrice,
			&item.Subtotal,
			&item.Notes,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal membaca order item: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("gagal membaca hasil order item: %w", err)
	}

	if items == nil {
		items = []OrderItemResponse{}
	}

	return items, nil
}

func generateInvoiceNo() string {
	return fmt.Sprintf(
		"INV-%s",
		time.Now().Format("20060102150405.000000000"),
	)
}
