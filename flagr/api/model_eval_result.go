/*
 * Flagr
 *
 * Flagr is a feature flagging, A/B testing and dynamic configuration microservice. The base path for all the APIs is \"/api/v1\". 
 *
 * API version: 1.1.8
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package api

type EvalResult struct {
	FlagID int64                   `json:"flagID,omitempty"`
	FlagKey string                 `json:"flagKey,omitempty"`
	FlagSnapshotID int64           `json:"flagSnapshotID,omitempty"`
	SegmentID int64                `json:"segmentID,omitempty"`
	VariantID int64                `json:"variantID,omitempty"`
	VariantKey string              `json:"variantKey,omitempty"`
	VariantAttachment *interface{} `json:"variantAttachment,omitempty"`
	EvalContext *EvalContext       `json:"evalContext,omitempty"`
	Timestamp string               `json:"timestamp,omitempty"`
	EvalDebugLog *EvalDebugLog     `json:"evalDebugLog,omitempty"`
}
