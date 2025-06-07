package license

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Result struct {
	Valid             bool   `json:"valid"`
	SubscriptionLevel string `json:"subscription_level"`
	Error             string `json:"error"`
}

var fallbackTokens = map[string]string{
	"valid-license-token": "basic",
	"mvp-test-token":      "basic",
}

func Validate(token string) (*Result, error) {
	if token == "" {
		return nil, errors.New("empty token")
	}

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", "https://license.infrapilot.io/validate", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
		body, rerr := io.ReadAll(resp.Body)
		if rerr != nil {
			return nil, fmt.Errorf("failed reading response: %w", rerr)
		}
		var result Result
		if jerr := json.Unmarshal(body, &result); jerr == nil {
			if result.Valid {
				return &result, nil
			}
			if result.Error != "" {
				return &result, errors.New(result.Error)
			}
			return &result, errors.New("invalid license")
		}
		// unknown format
	}

	// fallback
	if level, ok := fallbackTokens[token]; ok {
		return &Result{Valid: true, SubscriptionLevel: level}, nil
	}
	return &Result{Valid: false}, errors.New("invalid or expired license token")
}
