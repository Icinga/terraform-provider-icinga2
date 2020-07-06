package icinga2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/lrsmith/go-icinga2-api/iapi"
)

func resourceIcinga2Service() *schema.Resource {

	return &schema.Resource{
		Create: resourceIcinga2ServiceCreate,
		Read:   resourceIcinga2ServiceRead,
		Delete: resourceIcinga2ServiceDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ServiceName",
				ForceNew:    true,
			},
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Hostname",
				ForceNew:    true,
			},
			"check_command": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CheckCommand",
				ForceNew:    true,
			},
			"vars": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"templates": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Templates",
				ForceNew:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceIcinga2ServiceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*iapi.Server)

	hostname := d.Get("hostname").(string)
	name := d.Get("name").(string)
	checkcommand := d.Get("check_command").(string)

	vars := make(map[string]string)

	templates := make([]string, len(d.Get("templates").([]interface{})))
	for i, v := range d.Get("templates").([]interface{}) {
		templates[i] = v.(string)
	}

	// Normalize from map[string]interface{} to map[string]string
	iterator := d.Get("vars").(map[string]interface{})
	for key, value := range iterator {
		vars[key] = value.(string)
	}

	services, err := client.CreateService(name, hostname, checkcommand, vars, templates)
	if err != nil {
		return err
	}

	found := false
	for _, service := range services {
		if service.Name == hostname+"!"+name {
			d.SetId(hostname + "!" + name)
			found = true
		}
	}

	if !found {
		return fmt.Errorf("Failed to Create Service %s!%s : %s", hostname, name, err)
	}

	return nil

}

func resourceIcinga2ServiceRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*iapi.Server)

	hostname := d.Get("hostname").(string)
	name := d.Get("name").(string)

	services, err := client.GetService(name, hostname)
	if err != nil {
		return err
	}

	found := false
	for _, service := range services {
		if service.Name == hostname+"!"+name {
			d.SetId(hostname + "!" + name)
			d.Set("hostname", hostname)
			d.Set("check_command", service.Attrs.CheckCommand)
			d.Set("vars", service.Attrs.Vars)
			found = true
		}
	}

	if !found {
		return fmt.Errorf("Failed to Read Service %s!%s : %s", hostname, name, err)
	}

	return nil

}

func resourceIcinga2ServiceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*iapi.Server)

	hostname := d.Get("hostname").(string)
	name := d.Get("name").(string)

	err := client.DeleteService(name, hostname)
	if err != nil {
		return fmt.Errorf("Failed to Delete Service %s!%s : %s", hostname, name, err)
	}

	return nil
}
