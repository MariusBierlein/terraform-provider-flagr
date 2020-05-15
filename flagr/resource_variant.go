package flagr

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mariusbierlein/terraform-provider-flagr/flagr/api"
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
			"attachment": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "{}",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return jsonEqual(old, new)
				},
			},
		},
	}
}

func resourceVariantCreate(data *schema.ResourceData, i interface{}) error {
	a := i.(*api.APIClient)
	flagId := stringToInt64(data.Get("flag_id").(string))
	variant, _, _ := a.VariantApi.CreateVariant(context.TODO(), flagId, api.CreateVariantRequest{
		Key:        data.Get("key").(string),
		Attachment: expandAttachment(data.Get("attachment").(string)),
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
	data.Set("attachment", flattenAttachment(variant.Attachment))
	return nil
}

func resourceVariantUpdate(data *schema.ResourceData, i interface{}) error {
	a := i.(*api.APIClient)
	if data.HasChanges("key", "attachment") {
		flagId := stringToInt64(data.Get("flag_id").(string))
		variantId := stringToInt64(data.Id())
		//todo: error handling
		a.VariantApi.PutVariant(context.TODO(), flagId, variantId, api.PutVariantRequest{
			Key:        data.Get("key").(string),
			Attachment: expandAttachment(data.Get("attachment").(string)),
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
