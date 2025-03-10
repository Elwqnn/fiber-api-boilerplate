package repository

import (
	"context"
	"backend/pkg/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountRepository interface {
	Create(ctx context.Context, account *models.Account) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Account, error)
	FindByProviderID(ctx context.Context, provider, providerAccountID string) (*models.Account, error)
	Update(ctx context.Context, account *models.Account) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type accountRepository struct {
	db *gorm.DB
}
func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, account *models.Account) error {
	return r.db.WithContext(ctx).Create(account).Error
}

func (r *accountRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&accounts).Error
	return accounts, err
}

func (r *accountRepository) FindByProviderID(ctx context.Context, provider, providerAccountID string) (*models.Account, error) {
	var account models.Account
	err := r.db.WithContext(ctx).Where("provider = ? AND provider_account_id = ?", provider, providerAccountID).First(&account).Error
	return &account, err
}

func (r *accountRepository) Update(ctx context.Context, account *models.Account) error {
	return r.db.WithContext(ctx).Save(account).Error
}

func (r *accountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Account{}, id).Error
}
