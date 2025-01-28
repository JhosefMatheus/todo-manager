package authcontroller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo-manager/controllers/auth/dto"
	"todo-manager/models"
	"todo-manager/services/auth/responses"
	dbservice "todo-manager/services/db"

	"github.com/joho/godotenv"
)

func TestSignIn(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Errorf("Erro ao conectar no banco: %v", err)
	}

	db, err := dbservice.GetDbConnection()

	deleteSql := `
		delete from user;
	`

	if err != nil {
		t.Error(err)
	}

	sql := `
		insert into user (name, email, password)
		value ('Jhosef Matheus', 'jhosef.dev@gmail.com', sha2('9=0=y7MA5S>y', 256));
	`

	if _, err = db.Exec(sql); err != nil {
		t.Errorf("Erro ao inserir usuário: %v", err)
	}

	defer dbservice.CloseDbConnection(db)
	defer db.Exec(deleteSql)

	var insertedUser models.UserModel

	sql = `
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
			Message:      "Usuário autenticado com sucesso.",
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
