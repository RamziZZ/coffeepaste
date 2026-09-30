package categories

import (
	"context"
	"errors"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

// GET ALL
func (s *Service) GetAll(
	ctx context.Context,
) ([]Category, error) {
	return s.repository.FindAll(ctx)
}

// GET BY ID
func (s *Service) GetByID(
	ctx context.Context,
	id int,
) (*Category, error) {
	if id <= 0 {
		return nil, errors.New("id kategori tidak valid")
	}
	return s.repository.FindByID(ctx, id)
}

// CREATE
func (s *Service) Create(
	ctx context.Context,
	req CreateCategoryRequest,
) (*Category, error) {
	if req.Name == "" {
		return nil, errors.New("nama kategori wajib diisi")
	}
	return s.repository.Create(ctx, req)
}

// UPDATE
func (s *Service) Update(
	ctx context.Context,
	id int,
	req UpdateCategoryRequest,
) (*Category, error) {
	if id <= 0 {
		return nil, errors.New("id kategori tidak valid")
	}
	if req.Name == "" {
		return nil, errors.New("nama kategori wajib diisi")
	}

	return s.repository.Update(
		ctx,
		id,
		req,
	)
}

// DELETE
func (s *Service) Delete(
	ctx context.Context,
	id int,
) error {
	if id <= 0 {
		return errors.New("id kategori tidak valid")
	}

	return s.repository.Delete(ctx, id)
}
