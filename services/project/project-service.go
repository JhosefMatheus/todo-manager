package project_service

import (
	project_dto "todo-manager/controllers/project/dto"
	"todo-manager/models"
	project_responses "todo-manager/services/project/responses"
)

func Create(dto project_dto.CreateProjectDTO) (int, *models.BaseResponse, *project_responses.CreateProjectResponse) {
	return 0, nil, nil
}
