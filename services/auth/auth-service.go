package authservice

import (
	"todo-manager/models"
	"todo-manager/services/auth/responses"
)

func SignIn() responses.SignInResponse {
	response := responses.SignInResponse{
		BaseResponse: models.BaseResponse{
			Message:      "Usuário autenticado com sucesso.",
			AlertVariant: models.SuccessAlertVariant,
		},
	}

	return response
}
