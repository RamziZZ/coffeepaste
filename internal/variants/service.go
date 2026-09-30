package variants

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
) ([]ProductVariant, error) {
	return s.repository.FindAll(ctx)
}

func (s *Service) GetByID(
	ctx context.Context,
	id int64,
) (*ProductVariant, error) {
	if id <= 0 {
		return nil, errors.New("id variant tidak valid")
	}

	return s.repository.FindByID(ctx, id)
}

func (s *Service) Create(
	ctx context.Context,
	req CreateVariantRequest,
) (*ProductVariant, error) {
	req.Name = strings.TrimSpace(req.Name)

	if req.ProductID <= 0 {
		return nil, errors.New("product_id tidak valid")
	}

	if req.Name == "" {
		return nil, errors.New("nama variant wajib diisi")
	}

	return s.repository.Create(ctx, req)
}

func (s *Service) Update(
	ctx context.Context,
	id int64,
	req UpdateVariantRequest,
) (*ProductVariant, error) {
	req.Name = strings.TrimSpace(req.Name)

	if id <= 0 {
		return nil, errors.New("id variant tidak valid")
	}

	if req.ProductID <= 0 {
		return nil, errors.New("product_id tidak valid")
	}

	if req.Name == "" {
		return nil, errors.New("nama variant wajib diisi")
	}

	return s.repository.Update(ctx, id, req)
}

func (s *Service) Delete(
	ctx context.Context,
	id int64,
) error {
	if id <= 0 {
		return errors.New("id variant tidak valid")
	}

	return s.repository.Delete(ctx, id)
}
