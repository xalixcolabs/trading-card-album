package auth_application

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOAuthConfig *oauth2.Config
	once              sync.Once
	httpClient        = &http.Client{Timeout: 10 * time.Second}
)

type GoogleResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Verified bool   `json:"verified_email"`
	Picture  string `json:"picture"`
}

func ConfigGoogleAuth() *oauth2.Config {
	once.Do(func() {
		googleOAuthConfig = &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		}
	})
	return googleOAuthConfig
}

func GetGoogleResponse(token string) (*GoogleResponse, error) {
	reqURL, err := url.Parse("https://www.googleapis.com/oauth2/v1/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to parse userinfo URL: %w", err)
	}
	ptoken := fmt.Sprintf("Bearer %s", token)
	res := &http.Request{
		Method: "GET",
		URL:    reqURL,
		Header: map[string][]string{
			"Authorization": {ptoken},
		},
	}
	req, err := httpClient.Do(res)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch userinfo: %w", err)
	}
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read userinfo response: %w", err)
	}
	var data GoogleResponse
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse userinfo response: %w", err)
	}
	return &data, nil
}
