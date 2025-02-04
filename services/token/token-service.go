package token_service

import (
	"fmt"
	"os"
	"time"
	"todo-manager/models"
	"todo-manager/utils"

	"github.com/golang-jwt/jwt/v5"
)

func Generate(user models.UserModel) (token string, err error) {
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

func Verify(token string) (user *models.UserModel, err error) {
	tokenSecret := os.Getenv("TOKEN_SECRET")

	tokenData, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil })

	if err != nil {
		return nil, err
	}

	claims, ok := tokenData.Claims.(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	id, ok := claims["id"].(float64)

	if !ok {
		return nil, fmt.Errorf("Invalid id format in token.")
	}

	createdAtText, ok := claims["createdAt"].(string)

	if !ok {
		return nil, fmt.Errorf("Invalid createdAt format in token.")
	}

	createdAt, err := utils.TextToTime(createdAtText)

	if err != nil {
		return nil, err
	}

	updatedAtText, ok := claims["updatedAt"].(*string)

	var updatedAt *time.Time

	if updatedAtText != nil {
		if !ok {
			return nil, fmt.Errorf("Invalid updatedAt format in token.")
		}

		updatedAt, err = utils.TextToTime(*updatedAtText)

		if err != nil {
			return nil, err
		}
	} else {
		updatedAt = nil
	}

	user = &models.UserModel{
		Id:        int(id),
		Name:      claims["name"].(string),
		Email:     claims["email"].(string),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return user, nil
}
