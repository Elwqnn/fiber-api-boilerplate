package model

import (
	"time"

	"github.com/google/uuid"
)

// Account stores OAuth account information for users
// One user can have multiple accounts (e.g., sign in with Google AND GitHub)
type Account struct {
	BaseModel
	UserID uuid.UUID `json:"user_id" gorm:"not null;index"`
	Type   string    `json:"type" gorm:"not null;default:credentials" validate:"required,oneof=oauth credentials"`

	// Credentials-specific
	Password string `json:"password,omitempty" validate:"required_if=Type credentials,min=6,max=100"`

	// OAuth-specific
	Provider          string    `json:"provider,omitempty" validate:"required_if=Type oauth"`
	ProviderAccountID string    `json:"provider_account_id,omitempty" validate:"required_if=Type oauth"`
	RefreshToken      string    `json:"refresh_token,omitempty"`
	AccessToken       string    `json:"access_token,omitempty"`
	ExpiresAt         time.Time `json:"expires_at,omitempty"`
	TokenType         string    `json:"token_type,omitempty"`
	Scope             string    `json:"scope,omitempty"`
}
