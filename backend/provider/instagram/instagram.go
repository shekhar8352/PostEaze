package instagram

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/shekhar8352/PostEaze/utils"
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
	SubscribeToWebhooks(accessToken string, pageID string, fields []string) error
	GetPageDetails(accessToken string) (*PageDetailsResponse, error)
	GetMedia(accessToken string, igUserID string, after string) (*MediaResponse, error)
	GetMediaInsights(accessToken string, mediaID string, metrics []string) (*InsightsResponse, error)
	GetStoryInsights(accessToken string, mediaID string) (*InsightsResponse, error)
	GetProfileInsights(accessToken string, igUserID string, metrics []string, since int64, until int64) (*InsightsResponse, error)
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

type PageDetailsResponse struct {
	ID                string `json:"id"`
	Username          string `json:"username"`
	Name              string `json:"name"`
	Biography         string `json:"biography"`
	FollowersCount    int    `json:"followers_count"`
	FollowsCount      int    `json:"follows_count"`
	MediaCount        int    `json:"media_count"`
	ProfilePictureURL string `json:"profile_picture_url"`
	Website           string `json:"website"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// MediaResponse represents the response from Instagram media endpoint
type MediaResponse struct {
	Data   []MediaItem `json:"data"`
	Paging *Paging     `json:"paging,omitempty"`
}

type MediaItem struct {
	ID        string `json:"id"`
	MediaType string `json:"media_type"` // IMAGE, VIDEO, CAROUSEL_ALBUM
	MediaURL  string `json:"media_url,omitempty"`
	Caption   string `json:"caption,omitempty"`
	Timestamp string `json:"timestamp"`
	Permalink string `json:"permalink,omitempty"`
	IGUserID  string `json:"ig_id,omitempty"`
}

type Paging struct {
	Cursors *Cursors `json:"cursors,omitempty"`
	Next    string   `json:"next,omitempty"`
}

type Cursors struct {
	Before string `json:"before,omitempty"`
	After  string `json:"after,omitempty"`
}

// InsightsResponse represents the response from Instagram insights endpoint
type InsightsResponse struct {
	Data []InsightData `json:"data"`
}

type InsightData struct {
	Name        string         `json:"name"`
	Period      string         `json:"period,omitempty"`
	Values      []InsightValue `json:"values"`
	Title       string         `json:"title,omitempty"`
	Description string         `json:"description,omitempty"`
	ID          string         `json:"id,omitempty"`
}

type InsightValue struct {
	Value   interface{} `json:"value"` // Can be int or map
	EndTime string      `json:"end_time,omitempty"`
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

func (p *InstagramProviderImpl) SubscribeToWebhooks(accessToken string, igUserID string, fields []string) error {
	// Use the Instagram Graph API endpoint for subscribing to webhooks
	// igUserID is the Instagram Professional Account ID
	reqURL := fmt.Sprintf("https://graph.instagram.com/%s/subscribed_apps", igUserID)

	data := url.Values{}
	data.Set("access_token", accessToken)
	data.Set("subscribed_fields", strings.Join(fields, ","))

	utils.Logger.Info(context.Background(), "Subscribing to webhooks with URL: %s", reqURL)

	resp, err := http.Post(reqURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	utils.Logger.Info(context.Background(), "Subscribed tto webhook response: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]interface{}
		if err := json.Unmarshal(body, &errResp); err == nil {
			return fmt.Errorf("failed to subscribe to webhooks: %v", errResp)
		}
		return fmt.Errorf("failed to subscribe to webhooks: %s, body: %s", resp.Status, string(body))
	}

	var successResp map[string]bool
	if err := json.Unmarshal(body, &successResp); err != nil {
		return fmt.Errorf("failed to parse subscription response: %w", err)
	}

	if !successResp["success"] {
		return fmt.Errorf("subscription returned success=false")
	}

	return nil
}

func (p *InstagramProviderImpl) GetPageDetails(accessToken string) (*PageDetailsResponse, error) {
	// Instagram Graph API endpoint to get user profile information
	reqURL := fmt.Sprintf("https://graph.instagram.com/me?fields=id,username,name,biography,followers_count,follows_count,media_count,profile_picture_url,website&access_token=%s", accessToken)

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
		var errResp map[string]interface{}
		if err := json.Unmarshal(body, &errResp); err == nil {
			return nil, fmt.Errorf("failed to get page details: %v", errResp)
		}
		return nil, fmt.Errorf("failed to get page details: %s, body: %s", resp.Status, string(body))
	}

	var pageDetails PageDetailsResponse
	if err := json.Unmarshal(body, &pageDetails); err != nil {
		return nil, fmt.Errorf("failed to parse page details response: %w", err)
	}

	return &pageDetails, nil
}

// GetMedia fetches media (posts/stories) for an Instagram user
func (p *InstagramProviderImpl) GetMedia(accessToken string, igUserID string, after string) (*MediaResponse, error) {
	fields := "id,media_type,media_url,caption,timestamp,permalink"
	reqURL := fmt.Sprintf("https://graph.instagram.com/%s/media?fields=%s&access_token=%s", igUserID, fields, accessToken)

	if after != "" {
		reqURL += "&after=" + after
	}

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
		var errResp map[string]interface{}
		if err := json.Unmarshal(body, &errResp); err == nil {
			return nil, fmt.Errorf("failed to get media: %v", errResp)
		}
		return nil, fmt.Errorf("failed to get media: %s, body: %s", resp.Status, string(body))
	}

	var mediaResp MediaResponse
	if err := json.Unmarshal(body, &mediaResp); err != nil {
		return nil, fmt.Errorf("failed to parse media response: %w", err)
	}

	return &mediaResp, nil
}

// GetMediaInsights fetches insights for a specific media (post/reel)
func (p *InstagramProviderImpl) GetMediaInsights(accessToken string, mediaID string, metrics []string) (*InsightsResponse, error) {
	metricsStr := strings.Join(metrics, ",")
	reqURL := fmt.Sprintf("https://graph.instagram.com/%s/insights?metric=%s&access_token=%s", mediaID, metricsStr, accessToken)

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
		var errResp map[string]interface{}
		if err := json.Unmarshal(body, &errResp); err == nil {
			return nil, fmt.Errorf("failed to get media insights: %v", errResp)
		}
		return nil, fmt.Errorf("failed to get media insights: %s, body: %s", resp.Status, string(body))
	}

	var insightsResp InsightsResponse
	if err := json.Unmarshal(body, &insightsResp); err != nil {
		return nil, fmt.Errorf("failed to parse insights response: %w", err)
	}

	return &insightsResp, nil
}

// GetStoryInsights fetches insights for a story
func (p *InstagramProviderImpl) GetStoryInsights(accessToken string, mediaID string) (*InsightsResponse, error) {
	// Story-specific metrics
	metrics := []string{"impressions", "reach", "exits", "replies", "taps_forward", "taps_back"}
	return p.GetMediaInsights(accessToken, mediaID, metrics)
}

// GetProfileInsights fetches insights for an Instagram profile
func (p *InstagramProviderImpl) GetProfileInsights(accessToken string, igUserID string, metrics []string, since int64, until int64) (*InsightsResponse, error) {
	metricsStr := strings.Join(metrics, ",")
	reqURL := fmt.Sprintf("https://graph.instagram.com/%s/insights?metric=%s&period=day&access_token=%s", igUserID, metricsStr, accessToken)

	if since > 0 {
		reqURL += fmt.Sprintf("&since=%d", since)
	}
	if until > 0 {
		reqURL += fmt.Sprintf("&until=%d", until)
	}

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
		var errResp map[string]interface{}
		if err := json.Unmarshal(body, &errResp); err == nil {
			return nil, fmt.Errorf("failed to get profile insights: %v", errResp)
		}
		return nil, fmt.Errorf("failed to get profile insights: %s, body: %s", resp.Status, string(body))
	}

	var insightsResp InsightsResponse
	if err := json.Unmarshal(body, &insightsResp); err != nil {
		return nil, fmt.Errorf("failed to parse profile insights response: %w", err)
	}

	return &insightsResp, nil
}
