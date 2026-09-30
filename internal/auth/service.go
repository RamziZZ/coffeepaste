package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"backend_coffepaste/pkg/jwt"
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

// REGISTER
func (s *Service) Register(
	ctx context.Context,
	req RegisterRequest,
) (*RegisterResponse, error) {

	// Cek apakah email sudah digunakan
	_, err := s.repository.FindByEmail(
		ctx,
		req.Email,
	)

	if err == nil {
		return nil, errors.New("email sudah terdaftar")
	}

	// Kalau error-nya bukan "data tidak ditemukan",
	// berarti ada masalah database
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	// Hash password
	passwordHash, err := password.HashPassword(
		req.Password,
	)

	if err != nil {
		return nil, err
	}

	// Simpan user sebagai CUSTOMER
	user, err := s.repository.Create(
		ctx,
		req.Name,
		req.Email,
		passwordHash,
	)

	if err != nil {
		return nil, err
	}

	return &RegisterResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  "CUSTOMER",
	}, nil
}

// LOGIN
func (s *Service) Login(
	ctx context.Context,
	req LoginRequest,
) (*LoginResponse, error) {

	// Cari user berdasarkan email
	user, err := s.repository.FindByEmail(
		ctx,
		req.Email,
	)

	if err != nil {

		// Email tidak ditemukan
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New(
				"email atau password salah",
			)
		}

		// Error database lainnya
		return nil, err
	}

	// Cek apakah akun aktif
	if !user.IsActive {
		return nil, errors.New(
			"akun tidak aktif",
		)
	}

	// Cek password
	err = password.CheckPassword(
		req.Password,
		user.PasswordHash,
	)

	if err != nil {
		return nil, errors.New(
			"email atau password salah",
		)
	}

	// Generate JWT
	token, err := jwt.GenerateToken(
		user.ID,
		user.RoleID,
		user.Role,
	)

	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User: UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
		Token: token,
	}, nil
}
