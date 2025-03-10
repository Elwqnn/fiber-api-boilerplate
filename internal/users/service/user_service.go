package service

import (
	"context"
	"backend/internal/users/repository"
	"backend/pkg/models"

	"github.com/google/uuid"
)

type UserService interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Create(ctx context.Context, user *models.User) error {
	return s.userRepo.Create(ctx, user)
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}

func (s *userService) Update(ctx context.Context, user *models.User) error {
	return s.userRepo.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}
