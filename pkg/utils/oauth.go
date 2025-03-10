package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type UserInfo struct {
	// User fields
	Name  string `json:"name"`
	Email string `json:"email"`
	Image string `json:"picture,omitempty"`

	// Account fields
	ID                interface{} `json:"id"`
	Provider          string      `json:"provider"`
	ProviderAccountID string      `json:"provider_account_id"`
	AccessToken       string      `json:"access_token"`
	TokenType         string      `json:"token_type"`
	Scope             string      `json:"scope"`
}

// GenerateRandomState creates a random state string for OAuth
func GenerateRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func (u *UserInfo) GetProviderID() string {
	switch v := u.ID.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', 0, 64)
	case int:
		return strconv.Itoa(v)
	case json.Number:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func GetUserInfoFromGoogle(accessToken string) (*UserInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		ID            interface{} `json:"id"`
		Email         string      `json:"email"`
		Name          string      `json:"name"`
		Picture       string      `json:"picture"`
		VerifiedEmail bool        `json:"verified_email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed decoding user info: %v", err)
	}

	return &UserInfo{
		ID:                result.ID,
		Name:              result.Name,
		Email:             result.Email,
		Image:             result.Picture,
		Provider:          "google",
		ProviderAccountID: fmt.Sprint(result.ID),
		AccessToken:       accessToken,
	}, nil
}

func GetUserInfoFromDiscord(accessToken string) (*UserInfo, error) {
	req, err := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		ID            string `json:"id"`
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
		Avatar        string `json:"avatar"`
		Email         string `json:"email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed decoding user info: %v", err)
	}

	return &UserInfo{
		ID:                result.ID,
		Name:              result.Username, // + "#" + result.Discriminator,
		Email:             result.Email,
		Image:             fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", result.ID, result.Avatar),
		Provider:          "discord",
		ProviderAccountID: result.ID,
		AccessToken:       accessToken,
	}, nil
}
