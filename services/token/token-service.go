package tokenservice

import (
	"os"
	"time"
	"todo-manager/models"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user models.UserModel) (token string, err error) {
	iat := time.Now()
	exp := iat.Add(48 * time.Hour).Unix()

	tokenObj := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":        user.Id,
			"name":      user.Name,
			"email":     user.Email,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
			"exp":       exp,
			"iat":       iat,
		},
	)

	tokenSecret := os.Getenv("TOKEN_SECRET")

	token, err = tokenObj.SignedString([]byte(tokenSecret))

	if err != nil {
		return "", err
	}

	return token, nil
}
