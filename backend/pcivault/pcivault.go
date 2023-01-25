// Package pcivault is responsible for making API requests to PCI Vault.
//
// Typically, you would declare an interface to call and implement,
// but I'm going to call these methods directly for the sake of brevity.
package pcivault

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"io"
	"net/http"
	"os"
)

type CaptureEndpoint struct {
	URL    string `json:"url"`
	Secret string `json:"secret"`
}

type TokenData struct {
	Token     string `json:"token"`
	Reference string `json:"reference,omitempty"`
}

const (
	baseURL string = "https://api.pcivault.io/v1"
)

func basicAuth(req *http.Request) {
	auth := os.Getenv("PCI_BASIC_AUTH")
	encoded := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", encoded))
}

func GetCaptureEndpoint() (CaptureEndpoint, error) {
	log.Info("Creating capture endpoint")

	keyUser := os.Getenv("PCI_KEY")
	keyPassphrase := os.Getenv("PCI_PASSPHRASE")

	url := fmt.Sprintf("%s/capture?user=%s&passphrase=%s", baseURL, keyUser, keyPassphrase)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return CaptureEndpoint{}, fmt.Errorf("failed to make request: %w", err)
	}
	basicAuth(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return CaptureEndpoint{}, fmt.Errorf("failed to create capture endpoint: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CaptureEndpoint{}, fmt.Errorf("failed to read PCI Vault response: %w", err)
	}
	if resp.StatusCode >= http.StatusBadRequest {
		log.Error(string(body))
		return CaptureEndpoint{}, fmt.Errorf("failed to create capture endpoint (see logs)")
	}

	var ce CaptureEndpoint
	err = json.Unmarshal(body, &ce)
	if err != nil {
		return ce, fmt.Errorf("failed to unmarshal capture endpoint result: %w", err)
	}

	log.Infof("Endpoint created: %s", ce.URL)
	return ce, nil
}

func GetVaultTokens() ([]TokenData, error) {
	log.Info("Retrieving list of tokens")

	keyUser := os.Getenv("PCI_KEY")

	url := fmt.Sprintf("%s/vault?user=%s", baseURL, keyUser)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET /vault request: %w", err)
	}
	basicAuth(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get tokens: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read PCI Vault response: %w", err)
	}
	if resp.StatusCode >= http.StatusBadRequest {
		log.Error(string(body))
		return nil, fmt.Errorf("failed to get list of tokens (see logs)")
	}

	var listResp map[string][]TokenData
	err = json.Unmarshal(body, &listResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal list response: %w", err)
	}

	tokens := listResp[keyUser]
	log.Infof("Retrieved a list of %d tokens", len(tokens))
	return tokens, nil
}
