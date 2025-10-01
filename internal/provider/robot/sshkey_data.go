package robot

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/cipherbullet/terraform-provider-hetzner/internal/types"
)

func DataSourceSSHKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSSHKeyRead,
		Schema: map[string]*schema.Schema{
			"fingerprint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data": {
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

func dataSourceSSHKeyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(types.Client)

	fingerprint := d.Get("fingerprint").(string)

	respBody, err := client.DoRequest("GET", "/key/"+fingerprint, nil)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return fmt.Errorf("ssh key with fingerprint %s not found", fingerprint)
		}
		return err
	}

	var resp SSHKeyResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return fmt.Errorf("error parsing ssh key response: %w", err)
	}

	d.SetId(resp.Key.Fingerprint)
	_ = d.Set("name", resp.Key.Name)
	_ = d.Set("data", resp.Key.Data)
	_ = d.Set("type", resp.Key.Type)
	_ = d.Set("size", resp.Key.Size)
	_ = d.Set("created_at", resp.Key.CreatedAt)

	return nil
}
