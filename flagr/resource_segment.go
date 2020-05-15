package flagr

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mariusbierlein/terraform-provider-flagr/flagr/api"
)

func ResourceSegment() *schema.Resource {
	return &schema.Resource{
		Create: resourceSegmentCreate,
		Read:   resourceSegmentRead,
		Update: resourceSegmentUpdate,
		Delete: resourceSegmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: resourceSegmentCustomDiff,
		Schema: map[string]*schema.Schema{
			"flag_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rollout_percent": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"rank": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"distributions": {
				Type:     schema.TypeSet,
				Optional: true,
				Set:      distributionSetFunc,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"percent": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"variant_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"variant_key": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceSegmentCustomDiff(diff *schema.ResourceDiff, i interface{}) error {
	// the flagr api doesn't allow to remove a segments distributions,
	// so we have to re-create the segment if that happens
	if diff.HasChange("distributions") {
		oldDist, newDist := diff.GetChange("distributions")
		if oldDist.(*schema.Set).Len() > 0 && newDist.(*schema.Set).Len() == 0 {
			return diff.ForceNew("distributions")
		}
	}
	return nil
}

func resourceSegmentCreate(data *schema.ResourceData, i interface{}) error {
	a := i.(*api.APIClient)
	flagId := stringToInt64(data.Get("flag_id").(string))
	segment, _, err := a.SegmentApi.CreateSegment(context.TODO(), flagId, api.CreateSegmentRequest{
		Description:    data.Get("description").(string),
		RolloutPercent: expandInt64(data.Get("rollout_percent")),
	})
	if err != nil {
		return err
	}
	data.SetId(int64ToString(segment.Id))
	return resourceSegmentUpdate(data, i)
}

func resourceSegmentRead(data *schema.ResourceData, i interface{}) error {
	a := i.(*api.APIClient)
	flagId := stringToInt64(data.Get("flag_id").(string))
	segments, _, err := a.SegmentApi.FindSegments(context.TODO(), flagId)
	var segment api.Segment
	for _, v := range segments {
		if v.Id == stringToInt64(data.Id()) {
			segment = v
			break
		}
	}
	if err != nil {
		data.SetId("")
		return err
	}
	data.Set("description", segment.Description)
	data.Set("rollout_percent", flattenInt64(segment.RolloutPercent))
	data.Set("distributions", flattenDistributions(segment.Distributions))
	return nil
}

func resourceSegmentUpdate(data *schema.ResourceData, i interface{}) error {
	a := i.(*api.APIClient)
	flagId := stringToInt64(data.Get("flag_id").(string))
	segmentId := stringToInt64(data.Id())
	if data.HasChanges("description", "rollout_percent") {
		//todo: error handling
		a.SegmentApi.PutSegment(context.TODO(), flagId, segmentId, api.PutSegmentRequest{
			Description:    data.Get("description").(string),
			RolloutPercent: expandInt64(data.Get("rollout_percent")),
		})
	}
	if data.HasChange("distributions") {
		a.DistributionApi.PutDistributions(context.TODO(), flagId, segmentId, api.PutDistributionsRequest{Distributions: expandDistributions(data.Get("distributions"))})
	}
	return resourceVariantRead(data, i)
}

func resourceSegmentDelete(data *schema.ResourceData, i interface{}) error {
	a := i.(*api.APIClient)
	flagId := stringToInt64(data.Get("flag_id").(string))
	segmentId := stringToInt64(data.Id())
	a.SegmentApi.DeleteSegment(context.TODO(), flagId, segmentId)
	return nil
}

func distributionSetFunc(i interface{}) int {
	return int(stringToInt64(i.(map[string]interface{})["variant_id"].(string)))
}
