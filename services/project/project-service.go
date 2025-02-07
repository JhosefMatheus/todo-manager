package project_service

import (
	"net/http"
	project_constants "todo-manager/constants/project"
	project_dto "todo-manager/controllers/project/dto"
	"todo-manager/models"
	db_service "todo-manager/services/db"
	project_responses "todo-manager/services/project/responses"
)

func Create(dto project_dto.CreateProjectDTO) (int, *models.BaseResponse, *project_responses.CreateProjectResponse) {
	db, err := db_service.GetDbConnection()

	if err != nil {
		responseBody := models.BaseResponse{
			Message:      project_constants.CreateProjectDbConnectionErrorMessage,
			AlertVariant: models.ErrorAlertVariant,
		}

		return http.StatusInternalServerError, &responseBody, nil
	}

	defer db_service.CloseDbConnection(db)

	if dto.ParentProjectId != nil {
		sql := `
			select count(*)
			from project
			where id = ?;
		`

		rows, err := db.Query(sql, dto.ParentProjectId)

		if err != nil {
			responseBody := models.BaseResponse{
				Message:      project_constants.CreateProjectDbConnectionErrorMessage,
				AlertVariant: models.ErrorAlertVariant,
			}

			return http.StatusInternalServerError, &responseBody, nil
		}

		var projectByIdCount int

		if rows.Next() {
			if err := rows.Scan(&projectByIdCount); err != nil {
				responseBody := models.BaseResponse{
					Message:      project_constants.CreateProjectDbConnectionErrorMessage,
					AlertVariant: models.ErrorAlertVariant,
				}

				return http.StatusInternalServerError, &responseBody, nil
			}
		} else {
			responseBody := models.BaseResponse{
				Message:      project_constants.ProjectNotFoundMessage,
				AlertVariant: models.WarningAlertVariant,
			}

			return http.StatusNotFound, &responseBody, nil
		}

		if projectByIdCount == 0 {
			responseBody := models.BaseResponse{
				Message:      project_constants.ProjectNotFoundMessage,
				AlertVariant: models.WarningAlertVariant,
			}

			return http.StatusNotFound, &responseBody, nil
		}
	}

	sql := `
		insert into project (parent_project_id, name)
		value (?, ?);
	`

	if _, err := db.Exec(sql, dto.ParentProjectId, dto.Name); err != nil {
		responseBody := models.BaseResponse{
			Message:      project_constants.CreateProjectDbConnectionErrorMessage,
			AlertVariant: models.ErrorAlertVariant,
		}

		return http.StatusInternalServerError, &responseBody, nil
	}

	responseBody := project_responses.CreateProjectResponse{
		BaseResponse: models.BaseResponse{
			Message:      project_constants.CreateProjectSuccessMessage,
			AlertVariant: models.SuccessAlertVariant,
		},
	}

	return http.StatusOK, nil, &responseBody
}
