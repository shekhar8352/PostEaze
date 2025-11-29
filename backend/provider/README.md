# Provider Package

This package contains implementations for external service providers.

## Available Providers

- [Meta Provider](#meta-provider)

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
