package authcontroller

import (
	"net/http"
	"todo-manager/controllers/auth/dto"
	"todo-manager/models"
	authservice "todo-manager/services/auth"

	"github.com/gin-gonic/gin"
)

func SignIn(c *gin.Context) {
	var requestBody dto.SignInDTO

	if err := c.BindJSON(&requestBody); err != nil {
		var response = models.BaseResponse{
			Message:      "DTO inválido",
			AlertVariant: models.WarningAlertVariant,
		}

		c.IndentedJSON(http.StatusForbidden, response)

		return
	}

	email, password := requestBody.Email, requestBody.Password

	status, errResponse, signInResponse := authservice.SignIn(email, password)

	if errResponse != nil {
		c.IndentedJSON(status, errResponse)
	} else {
		c.IndentedJSON(status, signInResponse)
	}
}
