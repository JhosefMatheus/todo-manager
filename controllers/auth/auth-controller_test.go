package authcontroller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	auth_constants "todo-manager/constants/auth"
	"todo-manager/controllers/auth/dto"
	"todo-manager/models"
	"todo-manager/services/auth/responses"
	dbservice "todo-manager/services/db"

	"github.com/joho/godotenv"
)

func TestInvalidRequestMethod(t *testing.T) {
	setupEnv(t)

	dto := dto.SignInDTO{}

	var dtoBytes bytes.Buffer

	if err := json.NewEncoder(&dtoBytes).Encode(dto); err != nil {
		t.Errorf("Erro ao codificar dto: %v", err)
	}

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/auth/sing-in", &dtoBytes)

	SignIn(res, req)

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
	setupEnv(t)

	dto := dto.SignInDTO{}

	var dtoBytes bytes.Buffer

	if err := json.NewEncoder(&dtoBytes).Encode(dto); err != nil {
		t.Errorf("Erro ao codificar dto: %v", err)
	}

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/sing-in", &dtoBytes)

	SignIn(res, req)

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
	setupEnv(t)

	db, err := dbservice.GetDbConnection()

	if err != nil {
		t.Errorf("Erro ao criar a conexão com o banco de dados: %v", err)
	}

	setupDB(db, t)

	defer dbservice.CloseDbConnection(db)
	defer clearDB(db)

	dto := dto.SignInDTO{
		Email:    "jhosef.dev@gmail.com",
		Password: "teste",
	}

	var dtoBytes bytes.Buffer

	if err = json.NewEncoder(&dtoBytes).Encode(dto); err != nil {
		t.Errorf("Erro ao codificar dto: %v", err)
	}

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/sing-in", &dtoBytes)

	SignIn(res, req)

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
	setupEnv(t)

	db, err := dbservice.GetDbConnection()

	if err != nil {
		t.Errorf("Erro ao criar a conexão com o banco de dados: %v", err)
	}

	setupDB(db, t)

	defer dbservice.CloseDbConnection(db)
	defer clearDB(db)

	var insertedUser models.UserModel

	sql := `
		select
			id,
			name,
			email,
			created_at,
			updated_at
		from user
		where email = ?
		limit 1;
	`

	rows, err := db.Query(sql, "jhosef.dev@gmail.com")

	if err != nil {
		t.Error(err)
	}

	if rows.Next() {
		if err = rows.Scan(&insertedUser.Id, &insertedUser.Name, &insertedUser.Email, &insertedUser.CreatedAt, &insertedUser.UpdatedAt); err != nil {
			t.Errorf("Erro ao ler resposta do banco de dados: %v", err)
		}
	} else {
		t.Error("Nenhum usuário encontrado")
	}

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

	SignIn(res, req)

	if http.StatusOK != res.Code {
		t.Errorf("Código de status esperado: %d, recebeu: %d", http.StatusOK, res.Code)
	}

	var resBody responses.SignInResponse

	if err = json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		t.Errorf("Erro ao decodificar o corpo da resposta: %v", err)
	}

	expectedBody := responses.SignInResponse{
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

func setupDB(db *sql.DB, t *testing.T) {
	sql := `
		insert into user (name, email, password)
		value ('Jhosef Matheus', 'jhosef.dev@gmail.com', sha2('9=0=y7MA5S>y', 256));
	`

	if _, err := db.Exec(sql); err != nil {
		t.Errorf("Erro ao inserir usuário: %v", err)
	}
}

func clearDB(db *sql.DB) {
	sql := `
		delete from user;
	`

	db.Exec(sql)
}

func setupEnv(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Errorf("Erro ao conectar no banco: %v", err)
	}
}
