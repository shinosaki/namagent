package types

import "encoding/json"

type Meta struct {
	Status    int    `json:"status"`
	ErrorCode string `json:"errorCode"`
}

type APIResponse struct {
	Meta Meta            `json:"meta"`
	Data json.RawMessage `json:"data"`
}
