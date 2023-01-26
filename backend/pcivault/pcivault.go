// Package pcivault is responsible for making API requests to PCI Vault.
//
// Typically, you would declare an interface to call and implement,
// but I'm going to call these methods directly for the sake of brevity.
package pcivault

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"io"
	"net/http"
	"os"
	"strconv"
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

func PostToStripe(td TokenData) error {
	log.Info("Posting token to Stripe")

	keyUser := os.Getenv("PCI_KEY")
	keyPassphrase := os.Getenv("PCI_PASSPHRASE")

	stripeKey := os.Getenv("STRIPE_KEY")

	// This map is here to demonstrate the data structure. In production code, you would probably want to use a struct.
	dataM := map[string]any{
		"request": map[string]any{
			"method": "POST",
			"url":    "https://api.stripe.com/v1/sources",
			"headers": []map[string]string{
				{"Content-Type": "application/x-www-form-urlencoded"},
				{"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(stripeKey+":"))},
			},
			"body": `type=card&card[number]={{card_number}}&card[exp_month]={{expiry_month}}&card[exp_year]={{expiry_year}}`,
		},
		"webhook": map[string]any{
			// You may use this free beeceptor endpoint to test your integration,
			// but please be mindful that it is rate-limited to 50 requests per day.
			// Rather use your own beeceptor (or any similar service) endpoint.
			"url": "https://vault.free.beeceptor.com/example",
		},
	}

	dataB, err := json.Marshal(dataM)
	if err != nil {
		return fmt.Errorf("failed to marshal data for Stripe Request: %w", err)
	}

	url := fmt.Sprintf("%s/proxy/post?user=%s&passphrase=%s&token=%s", baseURL, keyUser, keyPassphrase, td.Token)
	if len(td.Reference) > 0 {
		url = fmt.Sprintf("%s&reference=%s", url, td.Reference)
	}
	debugMode, _ := strconv.ParseBool(os.Getenv("DEBUG_MODE"))
	if debugMode {
		// Only do this in debug mode with test cards, never in production
		url = fmt.Sprintf("%s&debug=true", url)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(dataB))
	if err != nil {
		return fmt.Errorf("failed to make POST /proxy request: %w", err)
	}
	basicAuth(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to post proxy: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read PCI Vault response: %w", err)
	}
	if resp.StatusCode >= http.StatusBadRequest {
		log.Error(string(body))
		return fmt.Errorf("failed to send proxy (see logs)")
	}

	var respBody map[string]any
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	// The populated_template field will only be processed if the debug flag has been set to true in the request
	log.Infof("Populated template: %v", respBody["populated_template"])

	return nil
}
