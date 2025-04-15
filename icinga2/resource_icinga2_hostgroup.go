package icinga2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lrsmith/go-icinga2-api/iapi"
)

func resourceIcinga2Hostgroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceIcinga2HostgroupCreate,
		Read:   resourceIcinga2HostgroupRead,
		Update: resourceIcinga2HostgroupUpdate,
		Delete: resourceIcinga2HostgroupDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the HostGroup",
				ForceNew:    true,
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of HostGroup",
			},
		},
	}
}

func resourceIcinga2HostgroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*iapi.Server)
	name := d.Get("name").(string)
	displayName := d.Get("display_name").(string)

	_, err := client.CreateHostgroup(name, displayName)
	if err != nil {
		return err
	}

	return resourceIcinga2HostgroupRead(d, meta)
}

func resourceIcinga2HostgroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*iapi.Server)
	name := d.Get("name").(string)

	hostgroups, err := client.GetHostgroup(name)
	if err != nil {
		return err
	}

	found := false
	for _, hostgroup := range hostgroups {
		if hostgroup.Name == name {
			d.SetId(name)
			_ = d.Set("display_name", hostgroup.Attrs.DisplayName)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("Failed to Read Hostgroup %s : %s", name, err)
	}

	return nil
}

func resourceIcinga2HostgroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*iapi.Server)
	if d.HasChange("display_name") {
		name := d.Get("name").(string)
		displayName := d.Get("display_name").(string)
		params := &iapi.HostgroupParams{
			DisplayName: displayName,
		}
		_, err := client.UpdateHostgroup(name, params)
		if err != nil {
			return err
		}
	}

	return resourceIcinga2HostgroupRead(d, meta)
}

func resourceIcinga2HostgroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*iapi.Server)
	name := d.Get("name").(string)

	if err := client.DeleteHostgroup(name); err != nil {
		return fmt.Errorf("Failed to Delete Hostgroup %s : %s", name, err)
	}

	return nil
}
