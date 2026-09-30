package shifts

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

type Service interface {
	GetAll(ctx context.Context) ([]Shift, error)
	GetByID(ctx context.Context, id int64) (*Shift, error)
	GetByUserID(ctx context.Context, userID int64) ([]Shift, error)
	GetOpenByUserID(ctx context.Context, userID int64) (*Shift, error)
	Open(ctx context.Context, userID int64, req OpenShiftRequest) (*Shift, error)
	Close(ctx context.Context, userID int64, id int64, req CloseShiftRequest) (*Shift, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) GetAll(ctx context.Context) ([]Shift, error) {
	shifts, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	if shifts == nil {
		shifts = []Shift{}
	}

	return shifts, nil
}

func (s *service) GetByID(ctx context.Context, id int64) (*Shift, error) {
	if id <= 0 {
		return nil, errors.New("id shift tidak valid")
	}

	shift, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if shift == nil {
		return nil, errors.New("shift tidak ditemukan")
	}

	return shift, nil
}

func (s *service) GetByUserID(
	ctx context.Context,
	userID int64,
) ([]Shift, error) {
	if userID <= 0 {
		return nil, errors.New("user id tidak valid")
	}

	shifts, err := s.repository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if shifts == nil {
		shifts = []Shift{}
	}

	return shifts, nil
}

func (s *service) GetOpenByUserID(
	ctx context.Context,
	userID int64,
) (*Shift, error) {
	if userID <= 0 {
		return nil, errors.New("user id tidak valid")
	}

	return s.repository.FindOpenByUserID(ctx, userID)
}

func (s *service) Open(
	ctx context.Context,
	userID int64,
	req OpenShiftRequest,
) (*Shift, error) {
	if userID <= 0 {
		return nil, errors.New("user id tidak valid")
	}

	if req.OpeningCash < 0 {
		return nil, errors.New("opening cash tidak boleh negatif")
	}

	activeShift, err := s.repository.FindOpenByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if activeShift != nil {
		return nil, errors.New("kasir masih memiliki shift yang aktif")
	}

	shift, err := s.repository.Open(
		ctx,
		userID,
		req.OpeningCash,
	)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka shift: %w", err)
	}

	return shift, nil
}

func (s *service) Close(
	ctx context.Context,
	userID int64,
	id int64,
	req CloseShiftRequest,
) (*Shift, error) {
	if userID <= 0 {
		return nil, errors.New("user id tidak valid")
	}

	if id <= 0 {
		return nil, errors.New("id shift tidak valid")
	}

	if req.ClosingCash < 0 {
		return nil, errors.New("closing cash tidak boleh negatif")
	}

	shift, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if shift == nil {
		return nil, errors.New("shift tidak ditemukan")
	}

	if shift.UserID != userID {
		return nil, errors.New("shift bukan milik kasir tersebut")
	}

	if strings.ToUpper(shift.Status) != "OPEN" {
		return nil, errors.New("shift sudah ditutup")
	}

	closedShift, err := s.repository.Close(
		ctx,
		id,
		req.ClosingCash,
	)
	if err != nil {
		return nil, fmt.Errorf("gagal menutup shift: %w", err)
	}

	if closedShift == nil {
		return nil, errors.New("shift tidak dapat ditutup")
	}

	return closedShift, nil
}
