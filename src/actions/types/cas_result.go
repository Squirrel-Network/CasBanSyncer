package types

import "time"

type CasResult struct {
	Banned bool `json:"ok"`
	Result struct {
		TimeAdded time.Time `json:"time_added"`
	} `json:"result"`
}
