package responses

import "todo-manager/models"

type SignInResponse struct {
	models.BaseResponse
	User models.UserModel `json:"user"`
}
