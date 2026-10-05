package model

type SuccessResponse[T any] struct {
	Status string `json:"status"`
	Data   T      `json:"data"`
}
