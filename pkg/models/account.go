package models

import (
	"time"

	"github.com/google/uuid"
)

// Account stores OAuth account information for users
// One user can have multiple accounts (e.g., sign in with Google AND Discord)
// No validate keyword is used here because the creation of accounts is handled by the service layer
type Account struct {
	BaseModel
	UserID uuid.UUID `json:"user_id" gorm:"not null;index"`
	User   User      `json:"-" gorm:"foreignKey:UserID"`
	Type   string    `json:"type" gorm:"not null;default:'credentials'"`

	// Credentials-specific
	Password string `json:"password,omitempty"`

	// OAuth-specific
	Provider          string    `json:"provider,omitempty"`
	ProviderAccountID string    `json:"provider_account_id,omitempty"`
	RefreshToken      string    `json:"refresh_token,omitempty"`
	AccessToken       string    `json:"access_token,omitempty"`
	ExpiresAt         time.Time `json:"expires_at,omitempty"`
	TokenType         string    `json:"token_type,omitempty"`
	Scope             string    `json:"scope,omitempty"`
}
