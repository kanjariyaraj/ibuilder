package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kanjariyaraj/Builder/internal/errors"
)

type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	Interval        int    `json:"interval"`
	ExpiresIn       int    `json:"expires_in"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	Error       string `json:"error"`
	ErrorDesc   string `json:"error_description"`
}

func requestDeviceCode(clientID string) (*DeviceCodeResponse, error) {
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("scope", "repo,workflow")

	req, err := http.NewRequest("POST", DeviceCodeURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.Wrap(errors.KindInternal, "failed to create device code request", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(errors.KindNetwork, "failed to request device code", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(errors.KindInternal, "failed to read device code response", err)
	}

	var deviceResp DeviceCodeResponse
	if err := json.Unmarshal(body, &deviceResp); err != nil {
		return nil, errors.Wrap(errors.KindInternal, "failed to parse device code response", err)
	}

	if deviceResp.DeviceCode == "" {
		return nil, errors.New(errors.KindNetwork, "unexpected device code response")
	}

	return &deviceResp, nil
}

func pollAccessToken(clientID, deviceCode string, interval int) (*AccessTokenResponse, error) {
	deadline := time.Now().Add(15 * time.Minute)

	for time.Now().Before(deadline) {
		time.Sleep(time.Duration(interval) * time.Second)

		data := url.Values{}
		data.Set("client_id", clientID)
		data.Set("device_code", deviceCode)
		data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")

		req, err := http.NewRequest("POST", AccessTokenURL, strings.NewReader(data.Encode()))
		if err != nil {
			return nil, errors.Wrap(errors.KindInternal, "failed to create token poll request", err)
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, errors.Wrap(errors.KindNetwork, "failed to poll access token", err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, errors.Wrap(errors.KindInternal, "failed to read token response", err)
		}

		var tokenResp AccessTokenResponse
		if err := json.Unmarshal(body, &tokenResp); err != nil {
			return nil, errors.Wrap(errors.KindInternal, "failed to parse token response", err)
		}

		switch tokenResp.Error {
		case "":
			if tokenResp.AccessToken != "" {
				return &tokenResp, nil
			}
		case "authorization_pending":
			continue
		case "slow_down":
			interval += 5
			continue
		case "expired_token":
			return nil, errors.New(errors.KindNetwork, "device code expired, please try again")
		case "access_denied":
			return nil, errors.New(errors.KindPermission, "authorization denied by user")
		default:
			return nil, errors.New(errors.KindNetwork, fmt.Sprintf("unexpected error: %s", tokenResp.Error))
		}
	}

	return nil, errors.New(errors.KindNetwork, "authorization timed out after 15 minutes")
}

func Authenticate(clientID string) (*TokenData, error) {
	if clientID == "" {
		clientID = DefaultClientID
	}

	deviceResp, err := requestDeviceCode(clientID)
	if err != nil {
		return nil, err
	}

	fmt.Println("╔══════════════════════════════════════════════╗")
	fmt.Println("║       GitHub Device Authentication          ║")
	fmt.Println("╠══════════════════════════════════════════════╣")
	fmt.Printf("║ 1. Open: %s\n", deviceResp.VerificationURI)
	fmt.Println("║                                              ║")
	fmt.Printf("║ 2. Enter code: %s\n", deviceResp.UserCode)
	fmt.Println("║                                              ║")
	fmt.Println("║ Browser should open automatically.          ║")
	fmt.Println("║ Waiting for authorization...                ║")
	fmt.Println("╚══════════════════════════════════════════════╝")

	openBrowser(deviceResp.VerificationURI)

	tokenResp, err := pollAccessToken(clientID, deviceResp.DeviceCode, deviceResp.Interval)
	if err != nil {
		return nil, err
	}

	token := &TokenData{
		AccessToken: tokenResp.AccessToken,
		TokenType:   tokenResp.TokenType,
		Scope:       tokenResp.Scope,
	}

	if err := SaveToken(token); err != nil {
		return nil, err
	}

	fmt.Println("✓ Authentication successful!")
	return token, nil
}

func AuthStatus() (bool, *TokenData, error) {
	token, err := LoadToken()
	if err != nil {
		return false, nil, err
	}
	if token == nil {
		return false, nil, nil
	}

	if err := ValidateToken(token.AccessToken); err != nil {
		return false, token, err
	}

	return true, token, nil
}

func Logout() error {
	return DeleteToken()
}
