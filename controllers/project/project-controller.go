package project_controller

import (
	"encoding/json"
	"io"
	"net/http"
	project_constants "todo-manager/constants/project"
	project_dto "todo-manager/controllers/project/dto"
	"todo-manager/models"
	project_service "todo-manager/services/project"
)

func Create(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		responesBody := models.BaseResponse{
			Message:      project_constants.CreateProjectMethodNotAllowed,
			AlertVariant: models.WarningAlertVariant,
		}

		w.WriteHeader(http.StatusMethodNotAllowed)

		if err := json.NewEncoder(w).Encode(responesBody); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	body, err := io.ReadAll(req.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	var dto project_dto.CreateProjectDTO

	if err = json.Unmarshal(body, &dto); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if dto.IsInvalid() {
		message := ""

		if dto.IsNameInvalid() {
			message += project_constants.CreateProjectDTOInvalidNameMessage
		}

		responseBody := models.BaseResponse{
			Message:      message,
			AlertVariant: models.WarningAlertVariant,
		}

		w.WriteHeader(http.StatusForbidden)

		if err := json.NewEncoder(w).Encode(responseBody); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	status, errorResponseBody, successResponseBody := project_service.Create(dto)

	w.WriteHeader(status)

	if errorResponseBody != nil {
		if err := json.NewEncoder(w).Encode(errorResponseBody); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		if err := json.NewEncoder(w).Encode(successResponseBody); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
