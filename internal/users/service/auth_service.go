package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"backend/internal/users/repository"
	"backend/pkg/config"
	"backend/pkg/models"
	"backend/pkg/utils"

	"golang.org/x/oauth2"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthService interface {
	Register(ctx context.Context, user *models.User, password string) error
	Login(ctx context.Context, email, password string) error
	GetOAuthRedirectURL(provider string) (string, error)
	HandleOAuthCallback(ctx context.Context, provider, code, state string) (*models.User, error)
}

type authService struct {
	userRepo       repository.UserRepository
	accountRepo    repository.AccountRepository
	oauthProviders *config.OAuthProviders
}

func NewAuthService(
	userRepo repository.UserRepository,
	accountRepo repository.AccountRepository,
) AuthService {
	return &authService{
		userRepo:       userRepo,
		accountRepo:    accountRepo,
		oauthProviders: config.LoadOAuthConfig(),
	}
}

func (s *authService) Register(ctx context.Context, user *models.User, password string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	account := &models.Account{
		Type:     "credentials",
		Password: hashedPassword,
	}
	user.Accounts = []models.Account{*account}
	return s.userRepo.Create(ctx, user)
}

func (s *authService) Login(ctx context.Context, email, password string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return ErrUserNotFound
	}

	// Find credentials account
	var credAccount models.Account
	for _, acc := range user.Accounts {
		if acc.Type == "credentials" {
			credAccount = acc
			break
		}
	}

	// Verify password
	if !utils.CheckPassword(password, credAccount.Password) {
		return ErrInvalidCredentials
	}

	return nil
}

func (s *authService) GetOAuthRedirectURL(provider string) (string, error) {
	var config oauth2.Config
	switch provider {
	case "google":
		config = s.oauthProviders.Google.ToOAuth2Config()
	case "discord":
		config = s.oauthProviders.Discord.ToOAuth2Config()
	default:
		return "", fmt.Errorf("unsupported provider: %s", provider)
	}

	state := utils.GenerateRandomState()
	return config.AuthCodeURL(state), nil
}

func (s *authService) HandleOAuthCallback(ctx context.Context, provider, code, state string) (*models.User, error) {
	if code == "" {
		return nil, fmt.Errorf("authorization code is missing")
	}

	if state == "" {
		return nil, fmt.Errorf("state parameter is missing")
	}

	// Exchange code for token
	token, err := s.exchangeCodeForToken(ctx, provider, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info from provider
	userInfo, err := s.getUserInfo(provider, token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Add token info to userInfo
	userInfo.AccessToken = token.AccessToken
	userInfo.TokenType = token.Type() // Use Type() method instead of direct access

	// Get scopes from oauth config
	var scopes []string
	switch provider {
	case "google":
		scopes = s.oauthProviders.Google.Scopes
	case "discord":
		scopes = s.oauthProviders.Discord.Scopes
	}
	userInfo.Scope = strings.Join(scopes, " ")

	// Find or create user
	user, err := s.findOrCreateUser(ctx, userInfo, provider)
	if err != nil {
		return nil, fmt.Errorf("failed to process user: %w", err)
	}

	return user, nil
}

func (s *authService) exchangeCodeForToken(ctx context.Context, provider, code string) (*oauth2.Token, error) {
	var config oauth2.Config
	switch provider {
	case "google":
		config = s.oauthProviders.Google.ToOAuth2Config()
	case "discord":
		config = s.oauthProviders.Discord.ToOAuth2Config()
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %v", err)
	}

	return token, nil
}

func (s *authService) getUserInfo(provider, accessToken string) (*utils.UserInfo, error) {
	switch provider {
	case "google":
		return utils.GetUserInfoFromGoogle(accessToken)
	case "discord":
		return utils.GetUserInfoFromDiscord(accessToken)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

func (s *authService) findOrCreateUser(ctx context.Context, userInfo *utils.UserInfo, provider string) (*models.User, error) {
	existingAccount, err := s.accountRepo.FindByProviderID(ctx, provider, fmt.Sprint(userInfo.ID))
	if err == nil {
		return s.userRepo.FindByID(ctx, existingAccount.UserID)
	}

	// Try to find user by email
	existingUser, err := s.userRepo.FindByEmail(ctx, userInfo.Email)
	if err == nil {
		// User exists, update fields
		existingUser.Name = userInfo.Name
		existingUser.Image = userInfo.Image

		if err := s.userRepo.Update(ctx, existingUser); err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}

		// Create new OAuth account
		account := &models.Account{
			UserID:            existingUser.ID,
			Type:              "oauth",
			Provider:          provider,
			ProviderAccountID: fmt.Sprint(userInfo.ID),
			AccessToken:       userInfo.AccessToken,
			TokenType:         userInfo.TokenType,
			Scope:             userInfo.Scope,
		}

		if err := s.accountRepo.Create(ctx, account); err != nil {
			return nil, fmt.Errorf("failed to create OAuth account: %w", err)
		}

		return existingUser, nil
	}

	// Create new user and account
	user := &models.User{
		Name:  userInfo.Name,
		Email: userInfo.Email,
		Image: userInfo.Image,
		Role:  "user",
	}

	// Create user first
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create OAuth account
	account := &models.Account{
		UserID:            user.ID,
		Type:              "oauth",
		Provider:          provider,
		ProviderAccountID: fmt.Sprint(userInfo.ID),
		AccessToken:       userInfo.AccessToken,
		TokenType:         userInfo.TokenType,
		Scope:             userInfo.Scope,
	}

	if err := s.accountRepo.Create(ctx, account); err != nil {
		// Rollback user creation on error
		_ = s.userRepo.Delete(ctx, user.ID)
		return nil, fmt.Errorf("failed to create OAuth account: %w", err)
	}

	return user, nil
}
