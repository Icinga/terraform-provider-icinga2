package icinga2

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/lrsmith/go-icinga2-api/iapi"
)

var (
	errInsecureSSL = errors.New("Requests are only allowed to use the HTTPS protocol so that traffic remains encrypted")
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ICINGA2_API_URL", nil),
				Description: "The address of the Icinga2 server.",
			},
			"api_user": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ICINGA2_API_USER", nil),
				Description: "The user to authenticate to the Icinga2 Server as.",
			},
			"api_password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ICINGA2_API_PASSWORD", nil),
				Description: "The password for authenticating to the Icinga2 server.",
			},
			"insecure_skip_tls_verify": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: EnvBoolDefaultFunc("ICINGA2_INSECURE_SKIP_TLS_VERIFY", false),
				Description: "Disable TLS verify when connecting to Icinga2 Server.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"icinga2_host":         resourceIcinga2Host(),
			"icinga2_hostgroup":    resourceIcinga2Hostgroup(),
			"icinga2_checkcommand": resourceIcinga2Checkcommand(),
			"icinga2_service":      resourceIcinga2Service(),
			"icinga2_user":         resourceIcinga2User(),
			"icinga2_notification": resourceIcinga2Notification(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	config, _ := iapi.New(
		d.Get("api_user").(string),
		d.Get("api_password").(string),
		d.Get("api_url").(string),
		d.Get("insecure_skip_tls_verify").(bool),
	)

	if err := validateURL(d.Get("api_url").(string)); err != nil {
		return nil, err
	}

	if err := config.Connect(); err != nil {
		return nil, err
	}

	return config, nil
}

func validateURL(urlString string) error {
	tokens, err := url.Parse(urlString)
	if err != nil {
		return err
	}

	if tokens.Scheme != "https" {
		return errInsecureSSL
	}

	if !strings.HasSuffix(tokens.Path, "/v1") {
		return fmt.Errorf("error : Invalid API version %s specified. Only v1 is currently supported", tokens.Path)
	}

	return nil
}

// EnvBoolDefaultFunc is a helper function that returns
func EnvBoolDefaultFunc(k string, dv interface{}) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		if v := os.Getenv(k); v == "true" {
			return true, nil
		}

		return false, nil
	}
}
