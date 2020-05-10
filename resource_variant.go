package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mariusbierlein/terraform-provider-flagr/api"
)

func ResourceVariant() *schema.Resource {
	return &schema.Resource{
		Create: resourceVariantCreate,
		Read:   resourceVariantRead,
		Update: resourceVariantUpdate,
		Delete: resourceVariantDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"flag_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			// todo: implement
			//"attachment": {
			//	Type:     schema.TypeString,
			//	Optional: true,
			//},
		},
	}
}

func resourceVariantCreate(data *schema.ResourceData, i interface{}) error {
	a := i.(*api.APIClient)
	flagId := stringToInt64(data.Get("flag_id").(string))
	variant, _, _ := a.VariantApi.CreateVariant(context.TODO(), flagId, api.CreateVariantRequest{
		Key: data.Get("key").(string),
	})
	data.SetId(int64ToString(variant.Id))
	return resourceVariantRead(data, i)
}

func resourceVariantRead(data *schema.ResourceData, i interface{}) error {
	a := i.(*api.APIClient)
	flagId := stringToInt64(data.Get("flag_id").(string))
	variants, _, _ := a.VariantApi.FindVariants(context.TODO(), flagId)
	var variant api.Variant
	for _, v := range variants {
		if v.Id == stringToInt64(data.Id()) {
			variant = v
			break
		}
	}
	data.Set("key", variant.Key)
	return nil
}

func resourceVariantUpdate(data *schema.ResourceData, i interface{}) error {
	a := i.(*api.APIClient)
	if data.HasChange("key") {
		flagId := stringToInt64(data.Get("flag_id").(string))
		variantId := stringToInt64(data.Id())
		//todo: error handling
		a.VariantApi.PutVariant(context.TODO(), flagId, variantId, api.PutVariantRequest{
			Key: data.Get("key").(string),
		})
	}
	return resourceVariantRead(data, i)
}

func resourceVariantDelete(data *schema.ResourceData, i interface{}) error {
	a := i.(*api.APIClient)
	flagId := stringToInt64(data.Get("flag_id").(string))
	variantId := stringToInt64(data.Id())

	a.VariantApi.DeleteVariant(context.TODO(), flagId, variantId)
	return nil
}
