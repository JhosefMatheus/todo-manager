package token_responses

import "todo-manager/models"

type TokenVerifyResponse struct {
	models.BaseResponse
	User models.UserModel `json:"user"`
}
