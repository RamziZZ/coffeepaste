package orders

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

type Service interface {
	GetAll(ctx context.Context) ([]Order, error)
	GetByID(ctx context.Context, id int64) (*OrderDetailResponse, error)
	GetByCustomerID(ctx context.Context, customerID int64) ([]Order, error)
	Create(ctx context.Context, customerID int64, req CreateOrderRequest) (*OrderDetailResponse, error)
	AssignCashier(ctx context.Context, id int64, req AssignCashierRequest) (*Order, error)
	UpdateStatus(ctx context.Context, id int64, status string) (*Order, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) GetAll(ctx context.Context) ([]Order, error) {
	orders, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	if orders == nil {
		orders = []Order{}
	}

	return orders, nil
}

func (s *service) GetByID(
	ctx context.Context,
	id int64,
) (*OrderDetailResponse, error) {
	if id <= 0 {
		return nil, errors.New("id order tidak valid")
	}

	order, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, errors.New("order tidak ditemukan")
	}

	return order, nil
}

func (s *service) GetByCustomerID(
	ctx context.Context,
	customerID int64,
) ([]Order, error) {
	if customerID <= 0 {
		return nil, errors.New("customer id tidak valid")
	}

	orders, err := s.repository.FindByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	if orders == nil {
		orders = []Order{}
	}

	return orders, nil
}

func (s *service) Create(
	ctx context.Context,
	customerID int64,
	req CreateOrderRequest,
) (*OrderDetailResponse, error) {
	if customerID <= 0 {
		return nil, errors.New("customer id tidak valid")
	}

	orderType := strings.ToUpper(strings.TrimSpace(req.OrderType))

	if orderType != "DINE_IN" && orderType != "TAKEAWAY" {
		return nil, errors.New(
			"order type harus DINE_IN atau TAKEAWAY",
		)
	}

	req.OrderType = orderType

	if req.Discount < 0 {
		return nil, errors.New("discount tidak boleh negatif")
	}

	if req.Tax < 0 {
		return nil, errors.New("tax tidak boleh negatif")
	}

	if len(req.Items) == 0 {
		return nil, errors.New("order harus memiliki minimal satu item")
	}

	if orderType == "DINE_IN" {
		if req.TableNumber == nil ||
			strings.TrimSpace(*req.TableNumber) == "" {
			return nil, errors.New(
				"table number wajib diisi untuk DINE_IN",
			)
		}

		tableNumber := strings.TrimSpace(*req.TableNumber)
		req.TableNumber = &tableNumber
	}

	if orderType == "TAKEAWAY" {
		req.TableNumber = nil
	}

	for i := range req.Items {
		if req.Items[i].ProductID <= 0 {
			return nil, fmt.Errorf(
				"product id pada item %d tidak valid",
				i+1,
			)
		}

		if req.Items[i].Quantity <= 0 {
			return nil, fmt.Errorf(
				"quantity pada item %d harus lebih dari 0",
				i+1,
			)
		}

		if req.Items[i].Notes != nil {
			notes := strings.TrimSpace(*req.Items[i].Notes)

			if notes == "" {
				req.Items[i].Notes = nil
			} else {
				req.Items[i].Notes = &notes
			}
		}
	}

	invoiceNo := generateInvoiceNo()

	order, err := s.repository.Create(
		ctx,
		customerID,
		invoiceNo,
		req,
	)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *service) AssignCashier(
	ctx context.Context,
	id int64,
	req AssignCashierRequest,
) (*Order, error) {
	if id <= 0 {
		return nil, errors.New("id order tidak valid")
	}

	if req.CashierID <= 0 {
		return nil, errors.New("cashier id tidak valid")
	}

	if req.ShiftID <= 0 {
		return nil, errors.New("shift id tidak valid")
	}

	order, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, errors.New("order tidak ditemukan")
	}

	if order.Order.Status == "COMPLETED" {
		return nil, errors.New(
			"order yang sudah selesai tidak dapat diubah",
		)
	}

	if order.Order.Status == "CANCELLED" {
		return nil, errors.New(
			"order yang dibatalkan tidak dapat diubah",
		)
	}

	result, err := s.repository.AssignCashier(
		ctx,
		id,
		req.CashierID,
		req.ShiftID,
	)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, errors.New(
			"order tidak dapat diassign",
		)
	}

	return result, nil
}

func (s *service) UpdateStatus(
	ctx context.Context,
	id int64,
	status string,
) (*Order, error) {
	if id <= 0 {
		return nil, errors.New("id order tidak valid")
	}

	status = strings.ToUpper(strings.TrimSpace(status))

	if !isValidStatus(status) {
		return nil, errors.New("status order tidak valid")
	}

	order, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, errors.New("order tidak ditemukan")
	}

	currentStatus := strings.ToUpper(
		strings.TrimSpace(order.Order.Status),
	)

	if !isValidStatusTransition(currentStatus, status) {
		return nil, fmt.Errorf(
			"perubahan status dari %s ke %s tidak diperbolehkan",
			currentStatus,
			status,
		)
	}

	result, err := s.repository.UpdateStatus(
		ctx,
		id,
		status,
	)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, errors.New(
			"status order tidak dapat diubah",
		)
	}

	return result, nil
}

func isValidStatus(status string) bool {
	switch status {
	case "PENDING",
		"PAID",
		"PROCESSING",
		"READY",
		"COMPLETED",
		"CANCELLED":
		return true
	default:
		return false
	}
}

func isValidStatusTransition(
	current string,
	next string,
) bool {
	if current == next {
		return true
	}

	switch current {
	case "PENDING":
		return next == "PAID" || next == "CANCELLED"

	case "PAID":
		return next == "PROCESSING" || next == "CANCELLED"

	case "PROCESSING":
		return next == "READY" || next == "CANCELLED"

	case "READY":
		return next == "COMPLETED"

	case "COMPLETED":
		return false

	case "CANCELLED":
		return false

	default:
		return false
	}
}