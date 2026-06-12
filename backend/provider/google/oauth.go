package google

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	TokenURL    = "https://oauth2.googleapis.com/token"
	AuthURL     = "https://accounts.google.com/o/oauth2/v2/auth"
	UserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"

	ScopeDriveReadonly   = "https://www.googleapis.com/auth/drive.readonly"
	ScopeYouTubeUpload   = "https://www.googleapis.com/auth/youtube.upload"
	ScopeYouTubeReadonly = "https://www.googleapis.com/auth/youtube.readonly"
	ScopeUserEmail       = "https://www.googleapis.com/auth/userinfo.email"
)

type OAuthProvider struct {
	ClientID     string
	ClientSecret string
	HTTPClient   *http.Client
}

func NewOAuthProvider() *OAuthProvider {
	return &OAuthProvider{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		HTTPClient:   &http.Client{Timeout: 30 * time.Second},
	}
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

type UserInfo struct {
	Email string `json:"email"`
	ID    string `json:"id"`
}

func (p *OAuthProvider) ExchangeCode(code, redirectURI string) (*TokenResponse, error) {
	if p.ClientID == "" || p.ClientSecret == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET must be set")
	}
	form := url.Values{}
	form.Set("code", code)
	form.Set("client_id", p.ClientID)
	form.Set("client_secret", p.ClientSecret)
	form.Set("redirect_uri", redirectURI)
	form.Set("grant_type", "authorization_code")

	req, err := http.NewRequest(http.MethodPost, TokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("token exchange failed (%d): %s", resp.StatusCode, string(body))
	}

	var out TokenResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (p *OAuthProvider) RefreshToken(refreshToken string) (*TokenResponse, error) {
	if p.ClientID == "" || p.ClientSecret == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET must be set")
	}
	form := url.Values{}
	form.Set("refresh_token", refreshToken)
	form.Set("client_id", p.ClientID)
	form.Set("client_secret", p.ClientSecret)
	form.Set("grant_type", "refresh_token")

	req, err := http.NewRequest(http.MethodPost, TokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("token refresh failed (%d): %s", resp.StatusCode, string(body))
	}

	var out TokenResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (p *OAuthProvider) GetUserInfo(accessToken string) (*UserInfo, error) {
	req, err := http.NewRequest(http.MethodGet, UserInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("userinfo failed (%d): %s", resp.StatusCode, string(body))
	}

	var out UserInfo
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func BuildAuthURL(clientID, redirectURI, scopes, state string) string {
	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("redirect_uri", redirectURI)
	params.Set("response_type", "code")
	params.Set("scope", scopes)
	params.Set("access_type", "offline")
	params.Set("prompt", "consent")
	if state != "" {
		params.Set("state", state)
	}
	return AuthURL + "?" + params.Encode()
}
