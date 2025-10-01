package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/cipherbullet/terraform-provider-hetzner/internal/provider/robot"
	"github.com/cipherbullet/terraform-provider-hetzner/internal/provider/cloud"
)

func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"robot": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
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
						"base_url": {
							Type:     schema.TypeString,
							Optional: true,
							Default: "https://robot-ws.your-server.de",
						},
					},
				},
			},
			"cloud": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"token": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"base_url": {
							Type:     schema.TypeString,
							Optional: true,
							Default: "https://api.hetzner.cloud/v1",
						},
					},
				},
			},
		},
		ResourcesMap: ResourcesMap(),
		DataSourcesMap: DataSourcesMap(),
	}
	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return ProviderConfigure(ctx, d)
	}

	return provider
}

func ResourcesMap() map[string]*schema.Resource {
	resourceMap, _ := ResourcesMapWithErrors()
	return resourceMap
}

func ResourcesMapWithErrors() (map[string]*schema.Resource, error) {
	return mergeResourceMaps(
		robotResources,
		cloudResources,
	)
}

func DataSourcesMap() map[string]*schema.Resource {
	dataResourceMap, _ := dataResourceMapWithErrors()
	return dataResourceMap
}

func dataResourceMapWithErrors() (map[string]*schema.Resource, error) {
	return mergeResourceMaps(
		robotDataResources,
		cloudDataResources,
	)
}

func ProviderConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	c := &Clients{}

	// Robot block
    if v, ok := d.GetOk("robot"); ok {
        robots := v.([]interface{})
        if len(robots) > 0 && robots[0] != nil {
            cfg := robots[0].(map[string]interface{})
            robotClient := robot.New(
                cfg["user"].(string),
                cfg["password"].(string),
                cfg["base_url"].(string),
            )
            c.Robot = robotClient
        }
    }

    // Cloud block
    if v, ok := d.GetOk("cloud"); ok {
        clouds := v.([]interface{})
        if len(clouds) > 0 && clouds[0] != nil {
            cfg := clouds[0].(map[string]interface{})
            cloudClient := cloud.New(
                cfg["token"].(string),
                cfg["base_url"].(string),
            )
            c.Cloud = cloudClient
            if c.Robot == nil {
                c.Robot = cloudClient
            }
        }
    }

	// If we still don't have a Robot client, return an error
	if c.Robot == nil && c.Cloud == nil {
		return nil, diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "No authentication provided",
				Detail:   "You must configure either robot { user/password } or cloud { token }",
			},
		}
	}

	return c, nil
}

func mergeResourceMaps(ms ...map[string]*schema.Resource) (map[string]*schema.Resource, error) {
	merged := make(map[string]*schema.Resource)
	duplicates := []string{}

	for _, m := range ms {
		for k, v := range m {
			if _, ok := merged[k]; ok {
				duplicates = append(duplicates, k)
			}

			merged[k] = v
		}
	}

	var err error
	if len(duplicates) > 0 {
		err = fmt.Errorf("saw duplicates in mergeResourceMaps: %v", duplicates)
	}

	return merged, err
}