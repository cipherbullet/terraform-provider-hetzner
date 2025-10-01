package cloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/cipherbullet/terraform-provider-hetzner/internal/types"
)

// Ensure the client implements the interface
var _ types.Client = (*Client)(nil)

type Client struct {
	Token   string
	BaseURL string
}

func New(token, baseURL string) *Client {
	return &Client{
		Token:   token,
		BaseURL: baseURL,
	}
}

// DoRequest makes an HTTP request to the Hetzner Cloud API
func (c *Client) DoRequest(method, path string, payload interface{}) ([]byte, error) {
	var body io.Reader
	if payload != nil {
		switch v := payload.(type) {
		case string:
			body = strings.NewReader(v)
		default:
			jsonData, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("error marshaling payload: %w", err)
			}
			body = bytes.NewBuffer(jsonData)
		}
	}

	req, err := http.NewRequest(method, c.BaseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Cloud API → Bearer token
	req.Header.Set("Authorization", "Bearer "+c.Token)
	if payload != nil && body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// DoFormRequest makes an HTTP request with form-encoded data to the Hetzner Cloud API
func (c *Client) DoFormRequest(method, path string, values map[string][]string) ([]byte, error) {
	formValues := url.Values(values)
	req, err := http.NewRequest(method, c.BaseURL+path, strings.NewReader(formValues.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
