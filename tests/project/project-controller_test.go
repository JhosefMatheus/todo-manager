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
	db_service "todo-manager/services/db"
	project_responses "todo-manager/services/project/responses"
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

func TestCreateProjectInvalidDTO(t *testing.T) {
	test_utils.SetupEnv(t, "../../.env")

	dto := project_dto.CreateProjectDTO{}

	var dtoBytes bytes.Buffer

	if err := json.NewEncoder(&dtoBytes).Encode(dto); err != nil {
		t.Errorf("Erro ao codificar dto: %v", err)
	}

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/project", &dtoBytes)

	project_controller.Create(res, req)

	if http.StatusForbidden != res.Code {
		t.Errorf("Código de status esperado: %d, recebeu: %d", http.StatusForbidden, res.Code)
	}

	expectedBody := models.BaseResponse{
		Message:      project_constants.CreateProjectDTOInvalidNameMessage,
		AlertVariant: models.WarningAlertVariant,
	}

	var resBody models.BaseResponse

	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		t.Errorf("Erro ao decodificar corpo da resposta: %v", err)
	}

	if expectedBody != resBody {
		t.Errorf("Corpo de resposta esperado: %v, recebeu %v", expectedBody, resBody)
	}
}

func TestCreateProjectParentProjectNotFound(t *testing.T) {
	test_utils.SetupEnv(t, "../../.env")

	parentProjectId := 1

	dto := project_dto.CreateProjectDTO{
		Name:            "Projetos Pessoais",
		ParentProjectId: &parentProjectId,
	}

	var dtoBytes bytes.Buffer

	if err := json.NewEncoder(&dtoBytes).Encode(dto); err != nil {
		t.Errorf("Erro ao codificar dto: %v", err)
	}

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/project", &dtoBytes)

	project_controller.Create(res, req)

	if http.StatusNotFound != res.Code {
		t.Errorf("Código de status esperado: %d, recebeu: %d", http.StatusNotFound, res.Code)
	}

	expectedBody := models.BaseResponse{
		Message:      project_constants.ProjectNotFoundMessage,
		AlertVariant: models.WarningAlertVariant,
	}

	var resBody models.BaseResponse

	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		t.Errorf("Erro ao decodificar corpo da resposta: %v", err)
	}

	if expectedBody != resBody {
		t.Errorf("Corpo de resposta esperado: %v, recebeu: %v", expectedBody, resBody)
	}
}

func TestCreateProjectWithNoParentProject(t *testing.T) {
	test_utils.SetupEnv(t, "../../.env")

	db, err := db_service.GetDbConnection()

	if err != nil {
		t.Errorf("Erro ao criar conexão com o banco de dados: %v", err)
	}

	defer db_service.CloseDbConnection(db)
	defer test_utils.ClearProjectTable(db)

	dto := project_dto.CreateProjectDTO{
		Name:            "Projetos Pessoais",
		ParentProjectId: nil,
	}

	var dtoBytes bytes.Buffer

	if err := json.NewEncoder(&dtoBytes).Encode(dto); err != nil {
		t.Errorf("Erro ao codificar dto: %v", err)
	}

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/project", &dtoBytes)

	project_controller.Create(res, req)

	if http.StatusOK != res.Code {
		t.Errorf("Código de status esperado: %d, recebeu: %d", http.StatusOK, res.Code)
	}

	expectedBody := project_responses.CreateProjectResponse{
		BaseResponse: models.BaseResponse{
			Message:      project_constants.CreateProjectSuccessMessage,
			AlertVariant: models.SuccessAlertVariant,
		},
	}

	var resBody project_responses.CreateProjectResponse

	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		t.Errorf("Erro ao decodificar corpo de resposta: %v", err)
	}

	if expectedBody != resBody {
		t.Errorf("Corpo de resposta esperado: %v, recebeu: %v", expectedBody, resBody)
	}
}

func TestCreateProjectWithParentProject(t *testing.T) {
	test_utils.SetupEnv(t, "../../.env")

	db, err := db_service.GetDbConnection()

	if err != nil {
		t.Errorf("Erro ao conectar com o banco de dados: %v", err)
	}

	parentProjectId := int(test_utils.SetupProjectTable(db, t))

	defer db_service.CloseDbConnection(db)
	defer test_utils.ClearProjectTable(db)

	dto := project_dto.CreateProjectDTO{
		Name:            "Projetos Pessoais",
		ParentProjectId: &parentProjectId,
	}

	var dtoBytes bytes.Buffer

	if err := json.NewEncoder(&dtoBytes).Encode(dto); err != nil {
		t.Errorf("Erro ao codificar dto: %v", err)
	}

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/project", &dtoBytes)

	project_controller.Create(res, req)

	if http.StatusOK != res.Code {
		t.Errorf("Código de status esperado: %d, recebeu %d", http.StatusOK, res.Code)
	}

	expectedBody := project_responses.CreateProjectResponse{
		BaseResponse: models.BaseResponse{
			Message:      project_constants.CreateProjectSuccessMessage,
			AlertVariant: models.SuccessAlertVariant,
		},
	}

	var resBody project_responses.CreateProjectResponse

	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		t.Errorf("Erro ao decodificar corpo da resposta: %v", err)
	}

	if expectedBody != resBody {
		t.Errorf("Corpo de resposta esperado: %v, recebeu: %v", expectedBody, resBody)
	}
}
