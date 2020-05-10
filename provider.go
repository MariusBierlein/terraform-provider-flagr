package main

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mariusbierlein/terraform-provider-flagr/api"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  80,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"flagr_flag":               ResourceFlag(),
			"flagr_variant":            ResourceVariant(),
			"flagr_segment":            ResourceSegment(),
			"flagr_segment_constraint": ResourceSegmentConstraint(),
		},
		ConfigureFunc: func(data *schema.ResourceData) (i interface{}, err error) {
			cfg := api.NewConfiguration()
			cfg.BasePath = fmt.Sprintf("http://%v:%v/api/v1", data.Get("host"), data.Get("port"))
			client := api.NewAPIClient(cfg)

			return client, nil
		},
	}
}
