package flagr

import (
	"github.com/checkr/goflagr"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"net/url"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"api_path": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/api/v1",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"flagr_flag":               ResourceFlag(),
			"flagr_variant":            ResourceVariant(),
			"flagr_segment":            ResourceSegment(),
			"flagr_segment_constraint": ResourceSegmentConstraint(),
		},
		ConfigureFunc: func(data *schema.ResourceData) (i interface{}, err error) {
			cfg := goflagr.NewConfiguration()
			base, _ := url.Parse(data.Get("api_host").(string))
			path, _ := url.Parse(data.Get("api_path").(string))
			cfg.BasePath = base.ResolveReference(path).String()
			client := goflagr.NewAPIClient(cfg)

			return client, nil
		},
	}
}
