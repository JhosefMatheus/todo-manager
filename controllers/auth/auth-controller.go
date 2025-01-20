package authcontroller

import (
	"fmt"
	"net/http"
	authservice "todo-manager/services/auth"

	"github.com/gin-gonic/gin"
)

func SignIn(c *gin.Context) {
	signInResponse := authservice.SignIn()

	fmt.Print(signInResponse)

	c.IndentedJSON(http.StatusOK, signInResponse)
}
