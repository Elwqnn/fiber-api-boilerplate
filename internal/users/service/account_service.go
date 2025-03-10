package service

import (
	"context"
	"backend/internal/users/repository"
	"backend/pkg/models"

	"github.com/google/uuid"
)

type AccountService interface {
	Create(ctx context.Context, account *models.Account) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.Account, error)
	GetByProviderID(ctx context.Context, provider, providerAccountID string) (*models.Account, error)
	Update(ctx context.Context, account *models.Account) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type accountService struct {
	accountRepo repository.AccountRepository
}

func NewAccountService(accountRepo repository.AccountRepository) AccountService {
	return &accountService{accountRepo: accountRepo}
}

func (s *accountService) Create(ctx context.Context, account *models.Account) error {
	return s.accountRepo.Create(ctx, account)
}

func (s *accountService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.Account, error) {
	return s.accountRepo.FindByUserID(ctx, userID)
}

func (s *accountService) GetByProviderID(ctx context.Context, provider, providerAccountID string) (*models.Account, error) {
	return s.accountRepo.FindByProviderID(ctx, provider, providerAccountID)
}

func (s *accountService) Update(ctx context.Context, account *models.Account) error {
	return s.accountRepo.Update(ctx, account)
}

func (s *accountService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.accountRepo.Delete(ctx, id)
}
