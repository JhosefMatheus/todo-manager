package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	project_constants "todo-manager/constants/project"
	project_controller "todo-manager/controllers/project"
	project_dto "todo-manager/controllers/project/dto"
	"todo-manager/models"
	test_utils "todo-manager/utils/test"
)

func TestCreateProjectInvalidMethod(t *testing.T) {
	test_utils.SetupEnv(t, "../../.env")

	dto := project_dto.CreateProjectDTO{}

	var dtoBytes bytes.Buffer

	if err := json.NewEncoder(&dtoBytes).Encode(dto); err != nil {
		t.Errorf("Erro ao codificar dto: %v", err)
	}

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/project", &dtoBytes)

	project_controller.Create(res, req)

	if http.StatusMethodNotAllowed != res.Code {
		t.Errorf("Código de status esperado: %d, recebeu: %d", http.StatusMethodNotAllowed, res.Code)
	}

	expectedBody := models.BaseResponse{
		Message:      project_constants.CreateProjectMethodNotAllowed,
		AlertVariant: models.WarningAlertVariant,
	}

	var resBody models.BaseResponse

	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		t.Errorf("Erro ao decodificar o corpo da resposta; %v", err)
	}

	if expectedBody != resBody {
		t.Errorf("Corpo de resposta esperado: %v, recebeu: %v", expectedBody, resBody)
	}
}
