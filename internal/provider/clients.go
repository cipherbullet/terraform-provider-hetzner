package provider

import (
	"github.com/cipherbullet/terraform-provider-hetzner/internal/types"
)

// Ensure Clients implements types.Client interface
var _ types.Client = (*Clients)(nil)

// Clients holds all API clients for Robot and Cloud
type Clients struct {
	Robot types.Client
	Cloud types.Client
}

// DoRequest forwards the request to the Robot client
func (c *Clients) DoRequest(method, path string, payload interface{}) ([]byte, error) {
	return c.Robot.DoRequest(method, path, payload)
}

// DoFormRequest forwards the form request to the Robot client
func (c *Clients) DoFormRequest(method, path string, values map[string][]string) ([]byte, error) {
	return c.Robot.DoFormRequest(method, path, values)
}
