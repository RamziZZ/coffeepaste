package products

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

// GET PRODUCT ALL
func (s *Service) GetAll(
	ctx context.Context,
) ([]Product, error) {
	return s.repository.FindAll(ctx)
}

// GET PRODUCT BY ID
func (s *Service) GetByID(
	ctx context.Context,
	id int64,
) (*Product, error) {
	if id <= 0 {
		return nil, errors.New("id produk tidak valid")
	}

	return s.repository.FindByID(ctx, id)
}

// CREATE PRODUCT
func (s *Service) Create(
	ctx context.Context,
	req CreateProductRequest,
) (*Product, error) {
	req.Name = strings.TrimSpace(req.Name)

	if req.CategoryID <= 0 {
		return nil, errors.New("category_id tidak valid")
	}

	if req.Name == "" {
		return nil, errors.New("nama produk wajib diisi")
	}

	if req.Price < 0 {
		return nil, errors.New("harga produk tidak boleh negatif")
	}

	if req.Stock < 0 {
		return nil, errors.New("stock produk tidak boleh negatif")
	}

	return s.repository.Create(ctx, req)
}

// UPDATE PRODUCT
func (s *Service) Update(
	ctx context.Context,
	id int64,
	req UpdateProductRequest,
) (*Product, error) {
	req.Name = strings.TrimSpace(req.Name)

	if id <= 0 {
		return nil, errors.New("id produk tidak valid")
	}

	if req.CategoryID <= 0 {
		return nil, errors.New("category_id tidak valid")
	}

	if req.Name == "" {
		return nil, errors.New("nama produk wajib diisi")
	}

	if req.Price < 0 {
		return nil, errors.New("harga produk tidak boleh negatif")
	}

	if req.Stock < 0 {
		return nil, errors.New("stock produk tidak boleh negatif")
	}

	return s.repository.Update(ctx, id, req)
}

// DELETE PRODUCT
func (s *Service) Delete(
	ctx context.Context,
	id int64,
) error {
	if id <= 0 {
		return errors.New("id produk tidak valid")
	}
	return s.repository.Delete(ctx, id)
}
