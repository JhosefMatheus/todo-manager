package models

type BaseResponse struct {
	Message      string       `json:"message"`
	AlertVariant AlertVariant `json:"alertVariant"`
}
