package auth_service

import (
	"net/http"
	auth_constants "todo-manager/constants/auth"
	"todo-manager/controllers/auth/dto"
	"todo-manager/models"
	"todo-manager/services/auth/responses"
	dbservice "todo-manager/services/db"
	tokenservice "todo-manager/services/token"
)

func SignIn(signInDTO dto.SignInDTO) (code int, errResponse *models.BaseResponse, response *responses.SignInResponse) {
	db, err := dbservice.GetDbConnection()

	if err != nil {
		return http.StatusInternalServerError, &models.BaseResponse{
			Message:      auth_constants.SignInDbConnectionErrorMessage,
			AlertVariant: models.ErrorAlertVariant,
		}, nil
	}

	defer dbservice.CloseDbConnection(db)

	email, password := signInDTO.Email, signInDTO.Password

	sql := `
		select
			id,
			name,
			email,
			created_at,
			updated_at
		from user
		where email = ? and password = sha2(?, 256)
		limit 1;
	`

	rows, err := db.Query(sql, email, password)

	if err != nil {
		return http.StatusInternalServerError, &models.BaseResponse{
			Message:      auth_constants.SignInDbConnectionErrorMessage,
			AlertVariant: models.ErrorAlertVariant,
		}, nil
	}

	var user models.UserModel

	if rows.Next() {
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return http.StatusInternalServerError, &models.BaseResponse{
				Message:      auth_constants.SignInDbConnectionErrorMessage,
				AlertVariant: models.ErrorAlertVariant,
			}, nil
		}
	} else {
		return http.StatusUnauthorized, &models.BaseResponse{
			Message:      auth_constants.SignInUnauthorizedMessage,
			AlertVariant: models.WarningAlertVariant,
		}, nil
	}

	token, err := tokenservice.Generate(user)

	if err != nil {
		return http.StatusInternalServerError, &models.BaseResponse{
			Message:      auth_constants.SignInGenerateTokenErrorMessage,
			AlertVariant: models.ErrorAlertVariant,
		}, nil
	}

	return http.StatusOK, nil, &responses.SignInResponse{
		BaseResponse: models.BaseResponse{
			Message:      auth_constants.SignInSuccessMessage,
			AlertVariant: models.SuccessAlertVariant,
		},
		User:  user,
		Token: token,
	}
}
