package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mariusbierlein/terraform-provider-flagr/api"
)

func ResourceSegmentConstraint() *schema.Resource {
	return &schema.Resource{
		Create: resourceSegmentConstraintCreate,
		Read:   resourceSegmentConstraintRead,
		Update: resourceSegmentConstraintUpdate,
		Delete: resourceSegmentConstraintDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"flag_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"segment_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"operator": {
				Type:     schema.TypeString,
				Required: true,
			},
			"property": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSegmentConstraintCreate(data *schema.ResourceData, i interface{}) error {
	a := i.(*api.APIClient)
	flagId := stringToInt64(data.Get("flag_id").(string))
	segmentId := stringToInt64(data.Get("segment_id").(string))
	constraint, _, _ := a.ConstraintApi.CreateConstraint(context.TODO(), flagId, segmentId, api.CreateConstraintRequest{
		Property: data.Get("property").(string),
		Operator: data.Get("operator").(string),
		Value:    data.Get("value").(string),
	})
	data.SetId(int64ToString(constraint.Id))
	return resourceSegmentConstraintRead(data, i)
}

func resourceSegmentConstraintRead(data *schema.ResourceData, i interface{}) error {
	a := i.(*api.APIClient)
	flagId := stringToInt64(data.Get("flag_id").(string))
	segmentId := stringToInt64(data.Get("segment_id").(string))
	constraints, _, _ := a.ConstraintApi.FindConstraints(context.TODO(), flagId, segmentId)
	var constraint api.Constraint
	for _, c := range constraints {
		if c.Id == stringToInt64(data.Id()) {
			constraint = c
			break
		}
	}
	data.Set("operator", constraint.Operator)
	data.Set("property", constraint.Property)
	data.Set("value", constraint.Value)
	return nil
}

func resourceSegmentConstraintUpdate(data *schema.ResourceData, i interface{}) error {
	a := i.(*api.APIClient)
	if data.HasChanges("operator", "property", "value") {
		flagId := stringToInt64(data.Get("flag_id").(string))
		segmentId := stringToInt64(data.Get("segment_id").(string))
		constraintId := stringToInt64(data.Id())
		//todo: error handling
		a.ConstraintApi.PutConstraint(context.TODO(), flagId, segmentId, constraintId, api.CreateConstraintRequest{
			Property: data.Get("property").(string),
			Operator: data.Get("operator").(string),
			Value:    data.Get("value").(string),
		})
	}
	return resourceSegmentConstraintRead(data, i)
}

func resourceSegmentConstraintDelete(data *schema.ResourceData, i interface{}) error {
	a := i.(*api.APIClient)
	flagId := stringToInt64(data.Get("flag_id").(string))
	segmentId := stringToInt64(data.Get("segment_id").(string))
	constraintId := stringToInt64(data.Id())

	a.ConstraintApi.DeleteConstraint(context.TODO(), flagId, segmentId, constraintId)
	return nil
}
