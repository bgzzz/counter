package model

//APIVersion current version of the API
const APIVersion = "v1"

//CounterRsp is a struct that is used for JSON response encoding
type CounterRsp struct {
	Counter uint64 `json:"counter"`
}
