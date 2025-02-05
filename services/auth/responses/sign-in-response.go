package auth_responses

import "todo-manager/models"

type SignInResponse struct {
	models.BaseResponse
	User  models.UserModel `json:"user"`
	Token string           `json:"token"`
}
