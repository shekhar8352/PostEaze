package mediasign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const defaultTTL = 6 * time.Hour

func signingKey() ([]byte, error) {
	keyStr := os.Getenv("ENCRYPTION_KEY")
	if keyStr == "" {
		return nil, fmt.Errorf("ENCRYPTION_KEY not set")
	}
	return []byte("stream:" + keyStr), nil
}

func Sign(versionID int64, exp time.Time) (string, error) {
	key, err := signingKey()
	if err != nil {
		return "", err
	}
	payload := fmt.Sprintf("%d|%d", versionID, exp.Unix())
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func Verify(versionID int64, expUnix int64, sig string) bool {
	if sig == "" || expUnix <= 0 {
		return false
	}
	if time.Now().Unix() > expUnix {
		return false
	}
	expected, err := Sign(versionID, time.Unix(expUnix, 0))
	if err != nil {
		return false
	}
	return hmac.Equal([]byte(expected), []byte(sig))
}

func PublicAPIBase() string {
	if base := strings.TrimSpace(os.Getenv("API_PUBLIC_BASE_URL")); base != "" {
		return strings.TrimRight(base, "/")
	}
	host := strings.TrimSpace(os.Getenv("API_HOST"))
	if host == "" {
		host = "localhost:8080"
	}
	scheme := "https"
	if strings.HasPrefix(host, "localhost") || strings.HasPrefix(host, "127.0.0.1") {
		scheme = "http"
	}
	return fmt.Sprintf("%s://%s/api", scheme, host)
}

func SignedStreamURL(versionID int64, ttl time.Duration) (string, error) {
	if ttl <= 0 {
		ttl = defaultTTL
	}
	exp := time.Now().Add(ttl)
	sig, err := Sign(versionID, exp)
	if err != nil {
		return "", err
	}
	u, _ := url.Parse(fmt.Sprintf("%s/v1/media/stream/%d", PublicAPIBase(), versionID))
	q := u.Query()
	q.Set("exp", strconv.FormatInt(exp.Unix(), 10))
	q.Set("sig", sig)
	u.RawQuery = q.Encode()
	return u.String(), nil
}
