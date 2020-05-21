package flagr

import (
	"context"
	"github.com/checkr/goflagr"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"strconv"
)

func ResourceFlag() *schema.Resource {
	return &schema.Resource{
		Create: resourceFlagCreate,
		Read:   resourceFlagRead,
		Update: resourceFlagUpdate,
		Delete: resourceFlagDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"data_records_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceFlagCreate(data *schema.ResourceData, i interface{}) error {
	a := i.(*goflagr.APIClient)
	flag, _, _ := a.FlagApi.CreateFlag(context.TODO(), goflagr.CreateFlagRequest{
		Description: data.Get("description").(string),
	})
	data.SetId(strconv.FormatInt(flag.Id, 10))
	//todo: only call update when needed to prevent unnecessary updates
	return resourceFlagUpdate(data, i)
}

func resourceFlagRead(data *schema.ResourceData, i interface{}) error {
	a := i.(*goflagr.APIClient)
	flag, _, _ := a.FlagApi.GetFlag(context.TODO(), stringToInt64(data.Id()))
	data.Set("key", flag.Key)
	data.Set("enabled", flag.Enabled)
	data.Set("description", flag.Description)
	data.Set("data_records_enabled", flag.DataRecordsEnabled)
	return nil
}

func resourceFlagUpdate(data *schema.ResourceData, i interface{}) error {
	a := i.(*goflagr.APIClient)
	flagId := stringToInt64(data.Id())
	if data.HasChanges("key", "description", "dataRecordsEnabled") {
		//todo: error handling
		a.FlagApi.PutFlag(context.TODO(), flagId, goflagr.PutFlagRequest{
			Description:        data.Get("description").(string),
			DataRecordsEnabled: data.Get("data_records_enabled").(bool),
			Key:                data.Get("key").(string),
		})
	}
	if data.HasChange("enabled") {
		// enabled in PutFlagRequest doesn't do anything, so this has to be a separate request
		a.FlagApi.SetFlagEnabled(context.TODO(), flagId, goflagr.SetFlagEnabledRequest{
			Enabled: data.Get("enabled").(bool),
		})
	}
	return resourceFlagRead(data, i)
}

func resourceFlagDelete(data *schema.ResourceData, i interface{}) error {
	a := i.(*goflagr.APIClient)
	id, _ := strconv.ParseInt(data.Id(), 10, 64)

	a.FlagApi.DeleteFlag(context.TODO(), id)
	return nil
}
