package auth

//this file implements the OAuth2 configuration and functions for Google authentication.
//It includes methods to generate the Google login URL and to retrieve user information after successful authentication.

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuthConfig struct {
	GoogleOAuthConfig *oauth2.Config
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func NewOAuthConfig(clientID, clientSecret, redirectURL string) *OAuthConfig {
	return &OAuthConfig{
		GoogleOAuthConfig: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (c *OAuthConfig) GetGoogleLoginURL(state string) string {
	return c.GoogleOAuthConfig.AuthCodeURL(state)
}

func (c *OAuthConfig) GetGoogleUserInfo(code string) (*GoogleUserInfo, error) {
	token, err := c.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()

	var userInfo GoogleUserInfo
	err = json.NewDecoder(response.Body).Decode(&userInfo)
	if err != nil {
		return nil, fmt.Errorf("failed decoding user info: %s", err.Error())
	}

	return &userInfo, nil
}
