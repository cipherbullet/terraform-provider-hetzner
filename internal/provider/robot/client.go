package robot

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

var _ types.Client = (*Client)(nil)

type Client struct {
	User     string
	Password string
	BaseURL  string
}

func New(user, password, baseURL string) *Client {
	return &Client{
		User:     user,
		Password: password,
		BaseURL:  baseURL,
	}
}

// DoRequest makes an HTTP request to the Hetzner Robot API
func (c *Client) DoRequest(method, path string, payload interface{}) ([]byte, error) {
	return c.doRequestWithContentType(method, path, payload, "application/json")
}

// DoFormRequest makes an HTTP request with form-encoded data to the Hetzner Robot API
func (c *Client) DoFormRequest(method, path string, values map[string][]string) ([]byte, error) {
	formValues := url.Values(values)
	return c.doRequestWithContentType(method, path, formValues.Encode(), "application/x-www-form-urlencoded")
}

func (c *Client) doRequestWithContentType(method, path string, payload interface{}, contentType string) ([]byte, error) {
	var body io.Reader

	switch v := payload.(type) {
	case nil:
		// no body
	case string:
		body = strings.NewReader(v)
	default:
		jsonData, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("error marshaling payload: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Robot API → Basic Auth
	req.SetBasicAuth(c.User, c.Password)
	if payload != nil {
		req.Header.Set("Content-Type", contentType)
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
