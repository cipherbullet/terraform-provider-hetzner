package robot

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/cipherbullet/terraform-provider-hetzner/internal/types"
)

type SSHKey struct {
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
	Type        string `json:"type"`
	Size        int    `json:"size"`
	Data        string `json:"data"`
	CreatedAt   string `json:"created_at"`
}

type SSHKeyResponse struct {
	Key SSHKey `json:"key"`
}

func ResourceSSHKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceSSHKeyCreate,
		Read:   resourceSSHKeyRead,
		Delete: resourceSSHKeyDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSSHKeyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(types.Client)

	values := url.Values{
		"name": {d.Get("name").(string)},
		"data": {d.Get("data").(string)},
	}

	respBody, err := client.DoFormRequest("POST", "/key", values)
	if err != nil {
		return err
	}

	var resp SSHKeyResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return fmt.Errorf("error parsing create response: %w", err)
	}

	// Use fingerprint as the Terraform resource ID
	d.SetId(resp.Key.Fingerprint)

	// Sync attributes
	_ = d.Set("name", resp.Key.Name)
	_ = d.Set("data", resp.Key.Data)
	_ = d.Set("fingerprint", resp.Key.Fingerprint)
	_ = d.Set("type", resp.Key.Type)
	_ = d.Set("size", resp.Key.Size)
	_ = d.Set("created_at", resp.Key.CreatedAt)

	return nil
}

func resourceSSHKeyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(types.Client)

	respBody, err := client.DoRequest("GET", "/key/"+d.Id(), nil)
	if err != nil {
		// If 404, resource no longer exists
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		}
		return err
	}

	var resp SSHKeyResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return fmt.Errorf("error parsing read response: %w", err)
	}

	_ = d.Set("name", resp.Key.Name)
	_ = d.Set("data", resp.Key.Data)
	_ = d.Set("fingerprint", resp.Key.Fingerprint)
	_ = d.Set("type", resp.Key.Type)
	_ = d.Set("size", resp.Key.Size)
	_ = d.Set("created_at", resp.Key.CreatedAt)

	return nil
}

func resourceSSHKeyDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(types.Client)

	_, err := client.DoRequest("DELETE", "/key/"+d.Id(), nil)
	return err
}
