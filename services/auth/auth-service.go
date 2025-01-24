package authservice

import (
	"fmt"
	"net/http"
	"todo-manager/models"
	"todo-manager/services/auth/responses"
	dbservice "todo-manager/services/db"
)

func SignIn(email string, password string) (code int, errResponse *models.BaseResponse, response *responses.SignInResponse) {
	db, err := dbservice.GetDbConnection()

	if err != nil {
		fmt.Print(err)

		return http.StatusInternalServerError, &models.BaseResponse{
			Message:      "Erro inesperado no banco de dados ao realizar o login.",
			AlertVariant: models.ErrorAlertVariant,
		}, nil
	}

	rows, err := db.Query("select id, name, email, created_at, updated_at from user where email = ? and password = sha2(?, 256) limit 1;", email, password)

	if err != nil {
		fmt.Print(err)

		return http.StatusInternalServerError, &models.BaseResponse{
			Message:      "Erro inesperado no banco de dados ao realizar o login.",
			AlertVariant: models.ErrorAlertVariant,
		}, nil
	}

	if !rows.Next() {
		return http.StatusUnauthorized, &models.BaseResponse{
			Message:      "Login ou senha inválido.",
			AlertVariant: models.WarningAlertVariant,
		}, nil
	}

	var user models.UserModel

	if rows.Next() {
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			fmt.Print(err)

			return http.StatusInternalServerError, &models.BaseResponse{
				Message:      "Erro inesperado no banco de dados ao realizar login.",
				AlertVariant: models.ErrorAlertVariant,
			}, nil
		}
	}

	fmt.Print(user)

	return http.StatusOK, nil, &responses.SignInResponse{
		BaseResponse: models.BaseResponse{
			Message:      "Usuário autenticado com sucesso.",
			AlertVariant: models.SuccessAlertVariant,
		},
		User: user,
	}
}
