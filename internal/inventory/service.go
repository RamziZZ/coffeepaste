package inventory

import (
	"context"
	"errors"
	"strings"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetAll(
	ctx context.Context,
) ([]InventoryMovement, error) {
	return s.repository.FindAll(ctx)
}

func (s *Service) GetByID(
	ctx context.Context,
	id int64,
) (*InventoryMovement, error) {
	if id <= 0 {
		return nil, errors.New("id inventory tidak valid")
	}

	return s.repository.FindByID(ctx, id)
}

func (s *Service) Create(
	ctx context.Context,
	req CreateInventoryMovementRequest,
) (*InventoryMovement, error) {
	req.Type = strings.ToUpper(strings.TrimSpace(req.Type))

	if req.ProductID <= 0 {
		return nil, errors.New("product_id tidak valid")
	}

	if req.Type == "" {
		return nil, errors.New("type inventory wajib diisi")
	}

	if req.Quantity <= 0 {
		return nil, errors.New("quantity harus lebih dari 0")
	}

	switch req.Type {
	case "IN", "OUT", "ADJUSTMENT":
	default:
		return nil, errors.New(
			"type inventory harus IN, OUT, atau ADJUSTMENT",
		)
	}

	if req.ReferenceID != nil {
		value := strings.TrimSpace(*req.ReferenceID)

		if value == "" {
			req.ReferenceID = nil
		} else {
			req.ReferenceID = &value
		}
	}

	if req.ReferenceType != nil {
		value := strings.TrimSpace(*req.ReferenceType)

		if value == "" {
			req.ReferenceType = nil
		} else {
			req.ReferenceType = &value
		}
	}

	if req.Note != nil {
		value := strings.TrimSpace(*req.Note)

		if value == "" {
			req.Note = nil
		} else {
			req.Note = &value
		}
	}

	return s.repository.Create(ctx, req)
}
