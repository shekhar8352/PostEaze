package instagram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	AuthURL           = "https://www.instagram.com/oauth/authorize"
	TokenURL          = "https://api.instagram.com/oauth/access_token"
	LongLivedTokenURL = "https://graph.instagram.com/access_token"
	RefreshTokenURL   = "https://graph.instagram.com/refresh_access_token"
)

type InstagramProvider interface {
	ExchangeCodeForToken(code string, redirectURI string) (*ShortLivedTokenResponse, error)
	GetLongLivedToken(shortLivedToken string) (*LongLivedTokenResponse, error)
	RefreshToken(accessToken string) (*LongLivedTokenResponse, error)
}

type InstagramProviderImpl struct {
	AppID     string
	AppSecret string
}

func NewInstagramProvider() *InstagramProviderImpl {
	return &InstagramProviderImpl{
		AppID:     os.Getenv("INSTAGRAM_APP_ID"),
		AppSecret: os.Getenv("INSTAGRAM_APP_SECRET"),
	}
}

type ShortLivedTokenResponse struct {
	AccessToken string `json:"access_token"`
	UserID      int64  `json:"user_id"`
}

type LongLivedTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"` // seconds until expiration (typically 5184000 = 60 days)
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func (p *InstagramProviderImpl) ExchangeCodeForToken(code string, redirectURI string) (*ShortLivedTokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", p.AppID)
	data.Set("client_secret", p.AppSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", redirectURI)
	data.Set("code", code)

	resp, err := http.Post(TokenURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to exchange code for token: %s, body: %s", resp.Status, string(body))
	}

	var tokenResp ShortLivedTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func (p *InstagramProviderImpl) GetLongLivedToken(shortLivedToken string) (*LongLivedTokenResponse, error) {
	reqURL := fmt.Sprintf("%s?grant_type=ig_exchange_token&client_secret=%s&access_token=%s",
		LongLivedTokenURL, p.AppSecret, shortLivedToken)

	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get long lived token: %s, body: %s", resp.Status, string(body))
	}

	var tokenResp LongLivedTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func (p *InstagramProviderImpl) RefreshToken(accessToken string) (*LongLivedTokenResponse, error) {
	reqURL := fmt.Sprintf("%s?grant_type=ig_refresh_token&access_token=%s",
		RefreshTokenURL, accessToken)

	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to refresh token: %s, body: %s", resp.Status, string(body))
	}

	var tokenResp LongLivedTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}
