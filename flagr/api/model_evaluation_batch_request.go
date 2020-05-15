/*
 * Flagr
 *
 * Flagr is a feature flagging, A/B testing and dynamic configuration microservice. The base path for all the APIs is \"/api/v1\".
 *
 * API version: 1.1.8
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package api

type EvaluationBatchRequest struct {
	Entities    []EvaluationEntity `json:"entities"`
	EnableDebug bool               `json:"enableDebug,omitempty"`
	// flagIDs
	FlagIDs []int64 `json:"flagIDs,omitempty"`
	// flagKeys. Either flagIDs or flagKeys works. If pass in both, Flagr may return duplicate results.
	FlagKeys []string `json:"flagKeys,omitempty"`
}
