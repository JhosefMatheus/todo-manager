package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	auth_constants "todo-manager/constants/auth"
	auth_controller "todo-manager/controllers/auth"
	auth_dto "todo-manager/controllers/auth/dto"
	"todo-manager/models"
	auth_responses "todo-manager/services/auth/responses"
	db_service "todo-manager/services/db"
	test_utils "todo-manager/utils/test"
)

func TestInvalidRequestMethod(t *testing.T) {
	test_utils.SetupEnv(t, "../../.env")

	dto := auth_dto.SignInDTO{}

	var dtoBytes bytes.Buffer

	if err := json.NewEncoder(&dtoBytes).Encode(dto); err != nil {
		t.Errorf("Erro ao codificar dto: %v", err)
	}

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/auth/sing-in", &dtoBytes)

	auth_controller.SignIn(res, req)

	if http.StatusMethodNotAllowed != res.Code {
		t.Errorf("Código de status esperado: %d, recebeu: %d", http.StatusMethodNotAllowed, res.Code)
	}

	expectedBody := models.BaseResponse{
		Message:      auth_constants.SignInMethodNotAllowedMessage,
		AlertVariant: models.WarningAlertVariant,
	}

	var resBody models.BaseResponse

	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		t.Errorf("Erro ao decodificar o corpo da resposta: %v", err)
	}

	if expectedBody != resBody {
		t.Errorf("Corpo de resposta esperado: %v, recebeu: %v", expectedBody, resBody)
	}
}

func TestInvalidDTO(t *testing.T) {
	test_utils.SetupEnv(t, "../../.env")

	dto := auth_dto.SignInDTO{}

	var dtoBytes bytes.Buffer

	if err := json.NewEncoder(&dtoBytes).Encode(dto); err != nil {
		t.Errorf("Erro ao codificar dto: %v", err)
	}

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/sing-in", &dtoBytes)

	auth_controller.SignIn(res, req)

	if http.StatusForbidden != res.Code {
		t.Errorf("Código de status esperado: %d, recebeu: %d", http.StatusForbidden, res.Code)
	}

	expectedBody := models.BaseResponse{
		Message:      fmt.Sprintf("%s\n%s", auth_constants.SignInInvalidEmailMessage, auth_constants.SignInInvalidPasswordMessage),
		AlertVariant: models.WarningAlertVariant,
	}

	var resBody models.BaseResponse

	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		t.Errorf("Erro ao decodificar o corpo da resposta: %v", err)
	}

	if expectedBody != resBody {
		t.Errorf("Corpo de resposta esperado: %v, corpo de resposta recebido: %v", expectedBody, resBody)
	}
}

func TestInvalidCredentials(t *testing.T) {
	test_utils.SetupEnv(t, "../../.env")

	db, err := db_service.GetDbConnection()

	if err != nil {
		t.Errorf("Erro ao criar a conexão com o banco de dados: %v", err)
	}

	test_utils.SetupUserTable(db, t)

	defer db_service.CloseDbConnection(db)
	defer test_utils.ClearUserTable(db)

	dto := auth_dto.SignInDTO{
		Email:    "jhosef.dev@gmail.com",
		Password: "teste",
	}

	var dtoBytes bytes.Buffer

	if err = json.NewEncoder(&dtoBytes).Encode(dto); err != nil {
		t.Errorf("Erro ao codificar dto: %v", err)
	}

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/sing-in", &dtoBytes)

	auth_controller.SignIn(res, req)

	if http.StatusUnauthorized != res.Code {
		t.Errorf("Código esperado: %d, código recebido: %d", http.StatusUnauthorized, res.Code)
	}

	expectedBody := models.BaseResponse{
		Message:      auth_constants.SignInUnauthorizedMessage,
		AlertVariant: models.WarningAlertVariant,
	}

	var resBody models.BaseResponse

	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		t.Errorf("Erro ao decodificar dto: %v", err)
	}

	if expectedBody != resBody {
		t.Errorf("Corpo de resposta esperado: %v, corpo recebido: %v", expectedBody, resBody)
	}
}

func TestSignIn(t *testing.T) {
	test_utils.SetupEnv(t, "../../.env")

	db, err := db_service.GetDbConnection()

	if err != nil {
		t.Errorf("Erro ao criar a conexão com o banco de dados: %v", err)
	}

	test_utils.SetupUserTable(db, t)

	defer db_service.CloseDbConnection(db)
	defer test_utils.ClearUserTable(db)

	insertedUser := test_utils.GetInsertedUser(db, t)

	dto := auth_dto.SignInDTO{
		Email:    "jhosef.dev@gmail.com",
		Password: "9=0=y7MA5S>y",
	}

	var dtoBytes bytes.Buffer

	if err = json.NewEncoder(&dtoBytes).Encode(dto); err != nil {
		t.Errorf("Erro ao codificar dto: %v", err)
	}

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/sign-in", &dtoBytes)

	auth_controller.SignIn(res, req)

	if http.StatusOK != res.Code {
		t.Errorf("Código de status esperado: %d, recebeu: %d", http.StatusOK, res.Code)
	}

	var resBody auth_responses.SignInResponse

	if err = json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		t.Errorf("Erro ao decodificar o corpo da resposta: %v", err)
	}

	expectedBody := auth_responses.SignInResponse{
		BaseResponse: models.BaseResponse{
			Message:      auth_constants.SignInSuccessMessage,
			AlertVariant: models.SuccessAlertVariant,
		},
		User: insertedUser,
	}

	if resBody.Message != expectedBody.Message {
		t.Errorf("Expected response message: %s, got: %s", expectedBody.Message, resBody.Message)
	}

	if resBody.AlertVariant != expectedBody.AlertVariant {
		t.Errorf("Expected alert variant: %s, got: %s", expectedBody.AlertVariant, resBody.AlertVariant)
	}

	if !resBody.User.Equals(expectedBody.User) {
		t.Errorf("Expected user: %v, got %v", expectedBody.User, resBody.User)
	}
}
