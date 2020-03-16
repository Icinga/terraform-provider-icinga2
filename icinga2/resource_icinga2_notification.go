package icinga2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/lrsmith/go-icinga2-api/iapi"
)

func resourceIcinga2Notification() *schema.Resource {

	return &schema.Resource{
		Create: resourceIcinga2NotificationCreate,
		Read:   resourceIcinga2NotificationRead,
		Delete: resourceIcinga2NotificationDelete,
		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"servicename": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"command": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"users": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"vars": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"templates": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"interval": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  1800,
			},
		},
	}
}

func resourceIcinga2NotificationCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*iapi.Server)

	hostname := d.Get("hostname").(string)
	servicename := d.Get("servicename").(string)
	command := d.Get("command").(string)
	interval := d.Get("interval").(int)
	var name string
	if servicename != "" {
		name = hostname + "!" + servicename + "!" + hostname + "-" + servicename
	} else {
		name = hostname + "!" + hostname
	}

	users := make([]string, len(d.Get("users").([]interface{})))
	for i, v := range d.Get("users").([]interface{}) {
		users[i] = v.(string)
	}

	vars := make(map[string]string)
	// Normalize from map[string]interface{} to map[string]string
	iterator := d.Get("vars").(map[string]interface{})
	for key, value := range iterator {
		vars[key] = value.(string)
	}

	templates := make([]string, len(d.Get("templates").([]interface{})))
	for i, v := range d.Get("templates").([]interface{}) {
		templates[i] = v.(string)
	}

	notifications, err := client.CreateNotification(name, hostname, command, servicename, interval, users, vars, templates)

	if err != nil {
		return err
	}

	found := false
	for _, notification := range notifications {
		if notification.Name == name {
			d.SetId(name)
			found = true
		}
	}

	if !found {
		return fmt.Errorf("Failed to Create Notification %s!%s : %s", hostname, servicename, err)
	}

	return nil
}

func resourceIcinga2NotificationRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*iapi.Server)

	hostname := d.Get("hostname").(string)
	servicename := d.Get("servicename").(string)
	var name string
	if servicename != "" {
		name = hostname + "!" + servicename + "!" + hostname + "-" + servicename
	} else {
		name = hostname + "!" + hostname
	}

	notifications, err := client.GetNotification(name)
	if err != nil {
		return err
	}

	found := false
	for _, notification := range notifications {
		if notification.Name == name {
			d.SetId(name)
			d.Set("hostname", hostname)
			d.Set("command", notification.Attrs.Command)
			d.Set("servicename", notification.Attrs.Servicename)
			d.Set("interval", notification.Attrs.Interval)
			d.Set("users", notification.Attrs.Users)
			d.Set("vars", notification.Attrs.Vars)
			found = true
		}
	}

	if !found {
		return fmt.Errorf("Failed to Read Notification %s!%s : %s", hostname, servicename, err)
	}

	return nil

}

func resourceIcinga2NotificationDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*iapi.Server)

	hostname := d.Get("hostname").(string)
	servicename := d.Get("servicename").(string)
	var name string
	if servicename != "" {
		name = hostname + "!" + servicename + "!" + hostname + "-" + servicename
	} else {
		name = hostname + "!" + hostname
	}

	err := client.DeleteNotification(name)
	if err != nil {
		return fmt.Errorf("Failed to Delete Notification %s!%s : %s", hostname, servicename, err)
	}

	return nil
}
