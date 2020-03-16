package icinga2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/lrsmith/go-icinga2-api/iapi"
)

func resourceIcinga2User() *schema.Resource {

	return &schema.Resource{
		Create: resourceIcinga2UserCreate,
		Read:   resourceIcinga2UserRead,
		Delete: resourceIcinga2UserDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username",
				ForceNew:    true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceIcinga2UserCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*iapi.Server)

	name := d.Get("name").(string)
	email := d.Get("email").(string)

	// Call CreateUser with normalized data
	users, err := client.CreateUser(name, email)
	if err != nil {
		return err
	}

	found := false
	for _, user := range users {
		if user.Name == name {
			d.SetId(name)
			found = true
		}
	}

	if !found {
		return fmt.Errorf("Failed to Create User %s : %s", name, err)
	}

	return nil
}

func resourceIcinga2UserRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*iapi.Server)

	name := d.Get("name").(string)

	users, err := client.GetUser(name)
	if err != nil {
		return err
	}

	found := false
	for _, user := range users {
		if user.Name == name {
			d.SetId(name)
			d.Set("name", user.Name)
			d.Set("email", user.Attrs.Email)
			found = true
		}
	}

	if !found {
		return fmt.Errorf("Failed to Read User %s : %s", name, err)
	}

	return nil
}

func resourceIcinga2UserDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*iapi.Server)
	name := d.Get("name").(string)

	err := client.DeleteUser(name)
	if err != nil {
		return fmt.Errorf("Failed to Delete User %s : %s", name, err)
	}

	return nil
}
