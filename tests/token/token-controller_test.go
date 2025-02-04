package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	token_constants "todo-manager/constants/token"
	auth_controller "todo-manager/controllers/auth"
	"todo-manager/controllers/auth/dto"
	token_controller "todo-manager/controllers/token"
	"todo-manager/models"
	authresponses "todo-manager/services/auth/responses"
	db_service "todo-manager/services/db"
	tokenresponses "todo-manager/services/token/responses"
	test_utils "todo-manager/utils/test"
)

func TestInvalidRequestMethod(t *testing.T) {
	test_utils.SetupEnv(t, "../../.env")

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/token/verify", nil)

	token_controller.Verify(res, req)

	if http.StatusMethodNotAllowed != res.Code {
		t.Errorf("Código de status esperado: %d, recebeu: %d", http.StatusMethodNotAllowed, res.Code)
	}

	expectedBody := models.BaseResponse{
		Message:      token_constants.TokenVerifyMethodNotAllowed,
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

func TestEmptyAuthorizationHeader(t *testing.T) {
	test_utils.SetupEnv(t, "../../.env")

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/token/verify", nil)

	req.Header.Set("Authorization", "")

	token_controller.Verify(res, req)

	if http.StatusUnauthorized != res.Code {
		t.Errorf("Código de status esperado: %d, recebeu: %d", http.StatusUnauthorized, res.Code)
	}

	expectedBody := models.BaseResponse{
		Message:      token_constants.TokenVerifyEmptyAuthorizationHeader,
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

func TestInvalidFormat(t *testing.T) {
	test_utils.SetupEnv(t, "../../.env")

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/token/verify", nil)

	req.Header.Set("Authorization", "token")

	token_controller.Verify(res, req)

	if http.StatusUnauthorized != res.Code {
		t.Errorf("Código de status esperado: %d, recebeu: %d", http.StatusUnauthorized, res.Code)
	}

	expectedBody := models.BaseResponse{
		Message:      token_constants.TokenVerifyInvalidFormat,
		AlertVariant: models.WarningAlertVariant,
	}

	var resBody models.BaseResponse

	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		t.Errorf("Erro ao decodificar corpo de resposta: %v", err)
	}

	if expectedBody != resBody {
		t.Errorf("Corpo de resposta esperado: %v, recebeu: %v", expectedBody, resBody)
	}
}

func TestInvalidToken(t *testing.T) {
	test_utils.SetupEnv(t, "../../.env")

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/token/verify", nil)

	req.Header.Set("Authorization", "Bearer token")

	token_controller.Verify(res, req)

	if http.StatusUnauthorized != res.Code {
		t.Errorf("Código de status esperado: %d, recebeu: %d", http.StatusUnauthorized, res.Code)
	}

	expectedBody := models.BaseResponse{
		Message:      token_constants.TokenVerifyInvalidToken,
		AlertVariant: models.WarningAlertVariant,
	}

	var resBody models.BaseResponse

	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		t.Errorf("Erro inesperado ao decodificar corpo da resposta: %v", err)
	}

	if expectedBody != resBody {
		t.Errorf("Corpo de resposta esperado: %v, recebeu: %v", expectedBody, resBody)
	}
}

func TestValidToken(t *testing.T) {
	test_utils.SetupEnv(t, "../../.env")

	db, err := db_service.GetDbConnection()

	if err != nil {
		t.Errorf("Erro ao criar a conexão com o banco de dados: %v", err)
	}

	test_utils.SetupUserTable(db, t)

	defer db_service.CloseDbConnection(db)
	defer test_utils.ClearUserTable(db)

	insertedUser := test_utils.GetInsertedUser(db, t)

	dto := dto.SignInDTO{
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

	var authSignInResBody authresponses.SignInResponse

	if err = json.NewDecoder(res.Body).Decode(&authSignInResBody); err != nil {
		t.Errorf("Erro ao decodificar o corpo da resposta: %v", err)
	}

	token := authSignInResBody.Token

	req = httptest.NewRequest(http.MethodGet, "/token/verify", nil)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	token_controller.Verify(res, req)

	var tokenVerifyResBody tokenresponses.TokenVerifyResponse

	if err = json.NewDecoder(res.Body).Decode(&tokenVerifyResBody); err != nil {
		t.Errorf("Erro ao decodificar o corpo da resposta: %v", err)
	}

	expectedBody := tokenresponses.TokenVerifyResponse{
		BaseResponse: models.BaseResponse{
			Message:      token_constants.TokenVerifySuccessMessage,
			AlertVariant: models.SuccessAlertVariant,
		},
		User: insertedUser,
	}

	if expectedBody.Message != tokenVerifyResBody.Message {
		t.Errorf("Mesagem de resposta esperada: %s, recebeu: %s", expectedBody.Message, tokenVerifyResBody.Message)
	}

	if expectedBody.AlertVariant != tokenVerifyResBody.AlertVariant {
		t.Errorf("Alert variant de resposta esperado: %s, recebeu: %s", expectedBody.AlertVariant, tokenVerifyResBody.AlertVariant)
	}

	if !expectedBody.User.Equals(tokenVerifyResBody.User) {
		t.Errorf("Usuário de resposta esperado: %v, recebeu: %v", expectedBody.User, tokenVerifyResBody.User)
	}
}
