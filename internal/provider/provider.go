package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cipherbullet/terraform-provider-hetzner/internal/cloud"
	"github.com/cipherbullet/terraform-provider-hetzner/internal/robot"
)

func New() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"user": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"robot_base_url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https://robot-ws.your-server.de",
			},
			"cloud_token": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"cloud_base_url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https://api.hetzner.cloud/v1",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			// Robot
			"hetzner_robot_boot":   robot.ResourceBoot(),
			"hetzner_robot_sshkey": robot.ResourceSSHKey(),

			// Cloud (future)
			// "hetzner_cloud_server": cloud.ResourceServer(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			// Robot
			"hetzner_robot_boot":   robot.DataSourceBoot(),
			"hetzner_robot_sshkey": robot.DataSourceSSHKey(),

			// Cloud (future)
			// "hetzner_cloud_server": cloud.DataSourceServer(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	c := &Clients{}

	// Initialize Robot client if credentials are provided
	if u, ok := d.GetOk("user"); ok {
		robotClient := robot.New(
			u.(string),
			d.Get("password").(string),
			d.Get("robot_base_url").(string),
		)
		c.Robot = robotClient
		c.Cloud = nil // We're using Robot, so Cloud can be nil
	}

	// Initialize Cloud client if token is provided
	if t, ok := d.GetOk("cloud_token"); ok {
		cloudClient := cloud.New(
			t.(string),
			d.Get("cloud_base_url").(string),
		)
		c.Cloud = cloudClient
		
		// If we don't have a Robot client, set it to the Cloud client
		// This allows the Clients struct to be used as a types.Client
		if c.Robot == nil {
			c.Robot = cloudClient
		}
	}

	// If we still don't have a Robot client, return an error
	if c.Robot == nil {
		return nil, fmt.Errorf("either user/password or cloud_token must be provided")
	}

	return c, nil
}
