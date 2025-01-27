package service

import (
	"context"
	"errors"
	"fiber-api-boilerplate/internal/config"
	"fiber-api-boilerplate/internal/model"
	"fiber-api-boilerplate/internal/repository"
	"fiber-api-boilerplate/pkg/utils"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthService interface {
	Register(ctx context.Context, user *model.User, password string) error
	Login(ctx context.Context, email, password string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	GetOAuthRedirectURL(provider string) (string, error)
	HandleOAuthCallback(ctx context.Context, provider, code, state string) (string, *model.User, error)
}

type authService struct {
	userRepo       repository.UserRepository
	accountRepo    repository.AccountRepository
	jwtSecret      string
	oauthProviders *config.OAuthProviders
}

func NewAuthService(
	userRepo repository.UserRepository,
	accountRepo repository.AccountRepository,
	jwtSecret string,
) AuthService {
	return &authService{
		userRepo:       userRepo,
		accountRepo:    accountRepo,
		jwtSecret:      jwtSecret,
		oauthProviders: config.LoadOAuthConfig(),
	}
}

func (s *authService) Register(ctx context.Context, user *model.User, password string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	account := &model.Account{
		Type:     "credentials",
		Password: hashedPassword,
	}
	user.Accounts = []model.Account{*account}
	return s.userRepo.Create(ctx, user)
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", ErrUserNotFound
	}

	// Find credentials account
	var credAccount model.Account
	for _, acc := range user.Accounts {
		if acc.Type == "credentials" {
			credAccount = acc
			break
		}
	}

	// Verify password
	if !utils.CheckPassword(password, credAccount.Password) {
		return "", ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID.String(), s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (s *authService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.jwtSecret), nil
	})
}

func (s *authService) GetOAuthRedirectURL(provider string) (string, error) {
	var config oauth2.Config
	switch provider {
	case "google":
		config = s.oauthProviders.Google.ToOAuth2Config()
	case "github":
		config = s.oauthProviders.Github.ToOAuth2Config()
	case "discord":
		config = s.oauthProviders.Discord.ToOAuth2Config()
	default:
		return "", fmt.Errorf("unsupported provider: %s", provider)
	}

	state := utils.GenerateRandomState()
	return config.AuthCodeURL(state), nil
}

func (s *authService) HandleOAuthCallback(ctx context.Context, provider, code, state string) (string, *model.User, error) {
	if code == "" {
		return "", nil, fmt.Errorf("authorization code is missing")
	}

	if state == "" {
		return "", nil, fmt.Errorf("state parameter is missing")
	}

	// Exchange code for token
	token, err := s.exchangeCodeForToken(ctx, provider, code)
	if err != nil {
		return "", nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info from provider
	userInfo, err := s.getUserInfo(ctx, provider, token.AccessToken)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Add token info to userInfo
	userInfo.AccessToken = token.AccessToken
	userInfo.TokenType = token.Type() // Use Type() method instead of direct access

	// Get scopes from oauth config
	var scopes []string
	switch provider {
	case "google":
		scopes = s.oauthProviders.Google.Scopes
	case "github":
		scopes = s.oauthProviders.Github.Scopes
	case "discord":
		scopes = s.oauthProviders.Discord.Scopes
	}
	userInfo.Scope = strings.Join(scopes, " ")

	// Find or create user
	user, err := s.findOrCreateUser(ctx, userInfo, provider)
	if err != nil {
		return "", nil, fmt.Errorf("failed to process user: %w", err)
	}

	// Generate JWT
	jwtToken, err := utils.GenerateToken(user.ID.String(), s.jwtSecret)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return jwtToken, user, nil
}

func (s *authService) exchangeCodeForToken(ctx context.Context, provider, code string) (*oauth2.Token, error) {
	var config oauth2.Config
	switch provider {
	case "google":
		config = s.oauthProviders.Google.ToOAuth2Config()
	case "github":
		config = s.oauthProviders.Github.ToOAuth2Config()
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

func (s *authService) getUserInfo(ctx context.Context, provider, accessToken string) (*utils.UserInfo, error) {
	switch provider {
	case "google":
		return utils.GetUserInfoFromGoogle(accessToken)
	case "github":
		return utils.GetUserInfoFromGithub(accessToken)
	case "discord":
		return utils.GetUserInfoFromDiscord(accessToken)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

func (s *authService) findOrCreateUser(ctx context.Context, userInfo *utils.UserInfo, provider string) (*model.User, error) {
	// First try to find existing account by provider ID
	existingAccount, err := s.accountRepo.FindByProviderID(ctx, provider, fmt.Sprint(userInfo.ID))
	if err == nil {
		// Account exists, get and return user
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
		account := &model.Account{
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
	user := &model.User{
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
	account := &model.Account{
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
