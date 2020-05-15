package main

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mariusbierlein/terraform-provider-flagr/api"
	"strconv"
)

func flattenDistributions(distributions []api.Distribution) *schema.Set {
	m := make([]interface{}, 0, 0)
	for _, d := range distributions {
		m = append(m, map[string]interface{}{
			"id":          int64ToString(d.Id),
			"percent":     flattenInt64(d.Percent),
			"variant_key": d.VariantKey,
			"variant_id":  int64ToString(d.VariantID),
		})
	}
	return schema.NewSet(distributionSetFunc, m)
}

func expandDistributions(list interface{}) []api.Distribution {
	distributions := make([]api.Distribution, 0, 0)
	for _, l := range list.(*schema.Set).List() {
		dis, ok := l.(map[string]interface{})
		if !ok {
			continue
		}
		distributions = append(distributions,
			api.Distribution{
				Id:         stringToInt64(dis["id"].(string)),
				Percent:    expandInt64(dis["percent"]),
				VariantKey: dis["variant_key"].(string),
				VariantID:  stringToInt64(dis["variant_id"].(string)),
			},
		)
	}
	return distributions
}

func flattenInt64(i int64) interface{} {
	return int(i)
}

func expandInt64(i interface{}) int64 {
	return int64(i.(int))
}

func stringToInt64(str string) int64 {
	out, _ := strconv.ParseInt(str, 10, 64)
	return out
}

func int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func expandAttachment(str string) *interface{} {
	var attachment interface{}
	json.Unmarshal([]byte(str), &attachment)
	return &attachment
}

func flattenAttachment(j interface{}) string {
	str, _ := json.Marshal(j)
	return string(str)
}
