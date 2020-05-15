/*
 * Flagr
 *
 * Flagr is a feature flagging, A/B testing and dynamic configuration microservice. The base path for all the APIs is \"/api/v1\".
 *
 * API version: 1.1.8
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package api

type PutFlagRequest struct {
	Description string `json:"description,omitempty"`
	// enabled data records will get data logging in the metrics pipeline, for example, kafka.
	DataRecordsEnabled bool `json:"dataRecordsEnabled,omitempty"`
	// it will overwrite entityType into evaluation logs if it's not empty
	EntityType string `json:"entityType,omitempty"`
	Enabled    bool   `json:"enabled,omitempty"`
	Key        string `json:"key,omitempty"`
	Notes      string `json:"notes,omitempty"`
}
