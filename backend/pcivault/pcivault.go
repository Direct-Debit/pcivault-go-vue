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
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
)

type CaptureEndpoint struct {
	URL    string `json:"url"`
	Secret string `json:"secret"`
}

const (
	baseURL string = "https://api.pcivault.io/v1"
)

func basicAuth() string {
	auth := os.Getenv("PCI_BASIC_AUTH")
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func GetCaptureEndpoint() (CaptureEndpoint, error) {
	keyUser := os.Getenv("PCI_KEY")
	keyPassphrase := os.Getenv("PCI_PASSPHRASE")

	url := fmt.Sprintf("%s/capture?user=%s&passphrase=%s", baseURL, keyUser, keyPassphrase)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return CaptureEndpoint{}, errors.Wrap(err, "failed to make request")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", basicAuth()))

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return CaptureEndpoint{}, errors.Wrap(err, "failed to create capture endpoint")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CaptureEndpoint{}, errors.Wrap(err, "failed to read PCI Vault response")
	}
	if resp.StatusCode >= http.StatusBadRequest {
		log.Error(string(body))
		return CaptureEndpoint{}, errors.New("failed to create capture endpoint (see logs)")
	}

	var ce CaptureEndpoint
	err = json.Unmarshal(body, &ce)
	return ce, errors.Wrap(err, "failed to unmarshal capture endpoint result")
}
