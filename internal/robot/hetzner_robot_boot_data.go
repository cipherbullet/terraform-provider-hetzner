package robot

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/cipherbullet/terraform-provider-hetzner/internal/types"
)

func DataSourceBoot() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBootRead,
		Schema: map[string]*schema.Schema{
			"server_number": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"keyboard": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssh_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBootRead(d *schema.ResourceData, m interface{}) error {
	client := m.(types.Client)

	server := d.Get("server_number").(string)
	mode := d.Get("mode").(string)

	if mode != "rescue" {
		return fmt.Errorf("only rescue mode supported in data source for now")
	}

	respBody, err := client.DoRequest("GET", "/boot/"+server+"/rescue", nil)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return fmt.Errorf("boot info for server %s not found", server)
		}
		return err
	}

	var cfg RescueConfig
	if err := json.Unmarshal(respBody, &cfg); err != nil {
		return fmt.Errorf("error parsing boot response: %w", err)
	}

	d.SetId(server)
	_ = d.Set("active", cfg.Rescue.Active)
	_ = d.Set("keyboard", cfg.Rescue.Keyboard)
	if len(cfg.Rescue.Authorized) > 0 {
		_ = d.Set("ssh_key", cfg.Rescue.Authorized[0])
	}

	return nil
}
