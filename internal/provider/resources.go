package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cipherbullet/terraform-provider-hetzner/internal/provider/robot"
	// "github.com/cipherbullet/terraform-provider-hetzner/internal/provider/cloud"
)

var robotResources = map[string]*schema.Resource {
	// Robot
	"hetzner_robot_boot":   robot.ResourceBoot(),
	"hetzner_robot_sshkey": robot.ResourceSSHKey(),
}

var robotDataResources = map[string]*schema.Resource {
	// Robot
	"hetzner_robot_boot":   robot.DataSourceBoot(),
	"hetzner_robot_sshkey": robot.DataSourceSSHKey(),
}

var cloudResources = map[string]*schema.Resource {
	// Cloud (future)
	// "hetzner_cloud_server": cloud.ResourceServer(),
}

var cloudDataResources = map[string]*schema.Resource {
	// Cloud (future)
	// "hetzner_cloud_server": cloud.DataSourceServer(),
}