package robot

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/cipherbullet/terraform-provider-hetzner/internal/types"
)

type RescueConfig struct {
	Rescue struct {
		ServerNumber string   `json:"server_number"`
		Active       bool     `json:"active"`
		Keyboard     string   `json:"keyboard"`
		Authorized   []string `json:"authorized_key"`
	} `json:"rescue"`
}

func ResourceBoot() *schema.Resource {
	return &schema.Resource{
		Create: resourceBootCreate,
		Read:   resourceBootRead,
		Delete: resourceBootDelete,

		Schema: map[string]*schema.Schema{
			"server_number": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"keyboard": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ssh_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceBootCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(types.Client)

	server := d.Get("server_number").(string)
	mode := d.Get("mode").(string)

	// Only rescue supported right now
	if mode != "rescue" {
		return fmt.Errorf("unsupported boot mode: %s (only 'rescue' is implemented for now)", mode)
	}

	// 1. Check if rescue already active
	respBody, err := client.DoRequest("GET", "/boot/"+server+"/rescue", nil)
	if err != nil {
		return err
	}

	var cfg RescueConfig
	if err := json.Unmarshal(respBody, &cfg); err != nil {
		return fmt.Errorf("error parsing rescue status: %w", err)
	}
	if cfg.Rescue.Active {
		return fmt.Errorf("rescue mode already active on server %s", server)
	}

	// 2. Activate rescue
	values := url.Values{}
	values.Set("keyboard", d.Get("keyboard").(string))
	values.Set("authorized_key", d.Get("ssh_key").(string))

	_, err = client.DoFormRequest("POST", "/boot/"+server+"/rescue", values)
	if err != nil {
		return err
	}

	// Use server_number as ID
	d.SetId(server)

	return resourceBootRead(d, m)
}

func resourceBootRead(d *schema.ResourceData, m interface{}) error {
	client := m.(types.Client)

	server := d.Get("server_number").(string)

	respBody, err := client.DoRequest("GET", "/boot/"+server+"/rescue", nil)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		}
		return err
	}

	var cfg RescueConfig
	if err := json.Unmarshal(respBody, &cfg); err != nil {
		return fmt.Errorf("error parsing rescue read: %w", err)
	}

	if !cfg.Rescue.Active {
		// rescue not active anymore → clear resource
		d.SetId("")
		return nil
	}

	// sync state
	_ = d.Set("server_number", cfg.Rescue.ServerNumber)
	_ = d.Set("keyboard", cfg.Rescue.Keyboard)
	if len(cfg.Rescue.Authorized) > 0 {
		_ = d.Set("ssh_key", cfg.Rescue.Authorized[0])
	}
	_ = d.Set("active", cfg.Rescue.Active)

	return nil
}

func resourceBootDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(types.Client)

	server := d.Get("server_number").(string)

	_, err := client.DoRequest("DELETE", "/boot/"+server+"/rescue", nil)
	if err != nil && !strings.Contains(err.Error(), "404") {
		return err
	}

	d.SetId("")
	return nil
}
