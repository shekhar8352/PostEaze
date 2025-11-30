package meta

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	GraphAPIURL = "https://graph.facebook.com/v18.0"
)

type MetaProvider interface {
	ExchangeCodeForToken(code string, redirectURI string) (*TokenResponse, error)
	GetLongLivedToken(shortLivedToken string) (*TokenResponse, error)
	GetPages(accessToken string) ([]Page, error)
}

type MetaProviderImpl struct {
	AppID     string
	AppSecret string
}

func NewMetaProvider() *MetaProviderImpl {
	return &MetaProviderImpl{
		AppID:     os.Getenv("META_APP_ID"),
		AppSecret: os.Getenv("META_APP_SECRET"),
	}
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type Page struct {
	ID                       string            `json:"id"`
	Name                     string            `json:"name"`
	AccessToken              string            `json:"access_token"`
	Category                 string            `json:"category"`
	Tasks                    []string          `json:"tasks"`
	InstagramBusinessAccount *InstagramAccount `json:"instagram_business_account,omitempty"`
}

type InstagramAccount struct {
	ID string `json:"id"`
}

type PagesResponse struct {
	Data []Page `json:"data"`
}

func (p *MetaProviderImpl) ExchangeCodeForToken(code string, redirectURI string) (*TokenResponse, error) {
	url := fmt.Sprintf("%s/oauth/access_token?client_id=%s&redirect_uri=%s&client_secret=%s&code=%s",
		GraphAPIURL, p.AppID, redirectURI, p.AppSecret, code)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to exchange code for token: %s", resp.Status)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func (p *MetaProviderImpl) GetLongLivedToken(shortLivedToken string) (*TokenResponse, error) {
	url := fmt.Sprintf("%s/oauth/access_token?grant_type=fb_exchange_token&client_id=%s&client_secret=%s&fb_exchange_token=%s",
		GraphAPIURL, p.AppID, p.AppSecret, shortLivedToken)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get long lived token: %s", resp.Status)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func (p *MetaProviderImpl) GetPages(accessToken string) ([]Page, error) {
	url := fmt.Sprintf("%s/me/accounts?access_token=%s&fields=id,name,access_token,category,tasks,instagram_business_account",
		GraphAPIURL, accessToken)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get pages: %s", resp.Status)
	}

	var pagesResp PagesResponse
	if err := json.NewDecoder(resp.Body).Decode(&pagesResp); err != nil {
		return nil, err
	}

	return pagesResp.Data, nil
}
