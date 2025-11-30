# Provider Package

This package contains implementations for external service providers.

## Available Providers

- [Meta Provider](#meta-provider)
- [Instagram Provider](#instagram-provider)

---

## Meta Provider

The `MetaProvider` handles interactions with the Meta (Facebook/Instagram) Graph API.

### Features
- Exchange authorization code for access token.
- Exchange short-lived access token for long-lived access token.
- Fetch Facebook Pages and connected Instagram Business Accounts.

### Usage

```go
import "github.com/shekhar8352/PostEaze/provider/meta"

provider := meta.NewMetaProvider()

// Exchange code
tokenResp, err := provider.ExchangeCodeForToken("auth_code", "redirect_uri")

// Get long-lived token
longLivedToken, err := provider.GetLongLivedToken(tokenResp.AccessToken)

// Get pages
pages, err := provider.GetPages(longLivedToken.AccessToken)
```

### Configuration
Requires the following environment variables:
- `META_APP_ID`
- `META_APP_SECRET`

---

## Instagram Provider

The `InstagramProvider` handles interactions with the Instagram Basic Display API for user authentication and channel creation.

### Features
- Exchange authorization code for short-lived access token
- Exchange short-lived token for long-lived access token (60 days)
- Refresh long-lived access tokens

### Usage

```go
import "github.com/shekhar8352/PostEaze/provider/instagram"

provider := instagram.NewInstagramProvider()

// Exchange code for short-lived token
shortToken, err := provider.ExchangeCodeForToken("auth_code", "redirect_uri")

// Get long-lived token (60 days)
longToken, err := provider.GetLongLivedToken(shortToken.AccessToken)

// Refresh token before expiration
refreshedToken, err := provider.RefreshToken(longToken.AccessToken)
```

### Configuration
Requires the following environment variables:
- `INSTAGRAM_APP_ID`
- `INSTAGRAM_APP_SECRET`
- `INSTAGRAM_REDIRECT_URI`
