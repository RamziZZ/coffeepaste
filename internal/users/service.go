package users

import (
	"context"
	"errors"
	"strings"

	"backend_coffepaste/pkg/password"
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
) ([]User, error) {
	return s.repository.FindAll(ctx)
}

func (s *Service) GetByID(
	ctx context.Context,
	id int64,
) (*User, error) {
	if id <= 0 {
		return nil, errors.New("id user tidak valid")
	}

	return s.repository.FindByID(ctx, id)
}

func (s *Service) CreateCashier(
	ctx context.Context,
	req CreateCashierRequest,
) (*User, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if req.Phone != nil {
		value := strings.TrimSpace(*req.Phone)

		if value == "" {
			req.Phone = nil
		} else {
			req.Phone = &value
		}
	}

	if req.Name == "" {
		return nil, errors.New("nama wajib diisi")
	}

	if req.Email == "" {
		return nil, errors.New("email wajib diisi")
	}

	if len(req.Password) < 6 {
		return nil, errors.New("password minimal 6 karakter")
	}

	existingUser, err := s.repository.FindByEmail(
		ctx,
		req.Email,
	)

	if err == nil && existingUser != nil {
		return nil, errors.New("email sudah digunakan")
	}

	passwordHash, err := password.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("gagal mengenkripsi password")
	}

	return s.repository.CreateCashier(
		ctx,
		req.Name,
		req.Email,
		passwordHash,
		req.Phone,
	)
}

func (s *Service) Update(
	ctx context.Context,
	id int64,
	req UpdateUserRequest,
) (*User, error) {
	if id <= 0 {
		return nil, errors.New("id user tidak valid")
	}

	if id == 1 {
		return nil, errors.New("akun owner utama tidak dapat diubah melalui endpoint ini")
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if req.Phone != nil {
		value := strings.TrimSpace(*req.Phone)

		if value == "" {
			req.Phone = nil
		} else {
			req.Phone = &value
		}
	}

	if req.Name == "" {
		return nil, errors.New("nama wajib diisi")
	}

	if req.Email == "" {
		return nil, errors.New("email wajib diisi")
	}

	existingUser, err := s.repository.FindByEmail(
		ctx,
		req.Email,
	)

	if err == nil && existingUser != nil && existingUser.ID != id {
		return nil, errors.New("email sudah digunakan")
	}

	user, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	if user.RoleID == 1 {
		return nil, errors.New("akun owner tidak dapat diubah")
	}

	return s.repository.Update(ctx, id, req)
}

func (s *Service) UpdatePassword(
	ctx context.Context,
	id int64,
	req UpdatePasswordRequest,
) error {
	if id <= 0 {
		return errors.New("id user tidak valid")
	}

	if id == 1 {
		return errors.New("password owner utama tidak diubah melalui endpoint users")
	}

	if len(req.Password) < 6 {
		return errors.New("password minimal 6 karakter")
	}

	user, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	if user.RoleID == 1 {
		return errors.New("password owner tidak dapat diubah melalui endpoint users")
	}

	passwordHash, err := password.HashPassword(req.Password)
	if err != nil {
		return errors.New("gagal mengenkripsi password")
	}

	return s.repository.UpdatePassword(
		ctx,
		id,
		passwordHash,
	)
}

func (s *Service) UpdateStatus(
	ctx context.Context,
	id int64,
	req UpdateStatusRequest,
) (*User, error) {
	if id <= 0 {
		return nil, errors.New("id user tidak valid")
	}

	if id == 1 {
		return nil, errors.New("status owner utama tidak dapat diubah")
	}

	user, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	if user.RoleID == 1 {
		return nil, errors.New("status owner tidak dapat diubah")
	}

	return s.repository.UpdateStatus(
		ctx,
		id,
		req.IsActive,
	)
}

func (s *Service) Delete(
	ctx context.Context,
	id int64,
) error {
	if id <= 0 {
		return errors.New("id user tidak valid")
	}

	if id == 1 {
		return errors.New("akun owner utama tidak dapat dihapus")
	}

	user, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	if user.RoleID == 1 {
		return errors.New("akun owner tidak dapat dihapus")
	}

	return s.repository.Delete(ctx, id)
}
