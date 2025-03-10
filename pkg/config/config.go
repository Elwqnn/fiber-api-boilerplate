package config

import (
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
	"golang.org/x/oauth2"
)

// Config is the main configuration struct
type Config struct {
	Port      string
	Env       string
	JWTSecret string
	Database  struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}
}

// OAuthConfig is the configuration struct for OAuth providers
type OAuthConfig struct {
	Provider     string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// OAuthProviders is the configuration struc containing all OAuth providers
type OAuthProviders struct {
	Google  OAuthConfig
	Discord OAuthConfig
}

var (
	GoogleEndpoints = oauth2.Endpoint{
		AuthURL:  "https://accounts.google.com/o/oauth2/auth",
		TokenURL: "https://oauth2.googleapis.com/token",
	}
	DiscordEnpoints = oauth2.Endpoint{
		AuthURL:  "https://discord.com/api/oauth2/authorize",
		TokenURL: "https://discord.com/api/oauth2/token",
	}
)

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Printf("Environment variable %s not set, using default value: %s", key, defaultValue)
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, strconv.Itoa(defaultValue))
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Error converting %s to integer, using default: %d", key, defaultValue)
		return defaultValue
	}
	return value
}

func LoadConfig() *Config {
	cfg := &Config{
		Port:      getEnv("PORT", "3000"),
		Env:       getEnv("ENV", "dev"),
		JWTSecret: getEnv("JWT_SECRET", "thisisaverylongsecret"),
	}

	cfg.Database.Host = getEnv("POSTGRES_HOST", "localhost")
	cfg.Database.Port = getEnv("POSTGRES_PORT", "5432")
	cfg.Database.User = getEnv("POSTGRES_USER", "postgres")
	cfg.Database.Password = getEnv("POSTGRES_PASSWORD", "postgres")
	cfg.Database.DBName = getEnv("POSTGRES_DB", "fiber-api-db")

	return cfg
}

func LoadOAuthConfig() *OAuthProviders {
	return &OAuthProviders{
		Google: OAuthConfig{
			Provider:     "google",
			ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:3000/api/v1/auth/callback/google"),
			Scopes:       []string{"profile", "email"},
		},
		Discord: OAuthConfig{
			Provider:     "discord",
			ClientID:     getEnv("DISCORD_CLIENT_ID", ""),
			ClientSecret: getEnv("DISCORD_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("DISCORD_REDIRECT_URL", "http://localhost:3000/api/v1/auth/callback/discord"),
			Scopes:       []string{"identify", "email"},
		},
	}
}

func (c *OAuthConfig) ToOAuth2Config() oauth2.Config {
	var endpoint oauth2.Endpoint

	switch c.Provider {
	case "google":
		endpoint = GoogleEndpoints
	case "discord":
		endpoint = DiscordEnpoints
	}

	return oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectURL,
		Scopes:       c.Scopes,
		Endpoint:     endpoint,
	}
}

func SetupSessionStore() *session.Store {
	storage := redis.New(redis.Config{
		Host:     getEnv("REDIS_HOST", "redis"),
		Port:     getEnvAsInt("REDIS_PORT", 6379),
		Password: getEnv("REDIS_PASSWORD", ""),
		Database: getEnvAsInt("REDIS_DB", 0),
	})

	sessionStore := session.New(session.Config{
		Storage:        storage,
		CookieSecure:   getEnv("ENV", "dev") == "prod", // HTTPS only
		CookieHTTPOnly: true,                           // Prevent client-side access
	})

	return sessionStore
}
