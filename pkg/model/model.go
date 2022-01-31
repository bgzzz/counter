package model

const APIVersion = "v1"

type CounterRsp struct {
	Counter uint64 `json:"counter"`
}
