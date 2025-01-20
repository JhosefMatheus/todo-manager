package authcontroller

import (
	"net/http"
	authservice "todo-manager/services/auth"

	"github.com/gin-gonic/gin"
)

func SignIn(c *gin.Context) {
	signInResponse := authservice.SignIn()

	c.IndentedJSON(http.StatusOK, signInResponse)
}
