package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/threpio/terraform-provider-blazemeter/internal/blazemeter"
)

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("RUNSCOPE_ACCESS_TOKEN", nil),
				Description: "A runscope access token.",
			},
			"api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RUNSCOPE_API_URL", nil),
				Description: "A runscope api url i.e. https://api.runscope.com.",
				Default:     "https://api.runscope.com",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{},

		ResourcesMap: map[string]*schema.Resource{},

		ConfigureContextFunc: providerConfigure,
	}
}

type providerConfig struct {
	client *blazemeter.Client
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	token := d.Get("access_token").(string)
	endpoint := d.Get("api_url").(string)

	client := blazemeter.NewClient(blazemeter.WithToken(token), blazemeter.WithEndpoint(endpoint))

	return &providerConfig{
		client: client,
	}, nil
}

func isNotFound(err error) bool {
	if blazemeterErr, ok := err.(blazemeter.Error); ok {
		return blazemeterErr.Status() == 404
	}
	return false
}