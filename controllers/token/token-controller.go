package token_controller

import (
	"encoding/json"
	"net/http"
	"strings"
	token_constants "todo-manager/constants/token"
	"todo-manager/models"
	token_service "todo-manager/services/token"
	"todo-manager/services/token/responses"
)

func Verify(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		responseBody := models.BaseResponse{
			Message:      token_constants.TokenVerifyMethodNotAllowed,
			AlertVariant: models.WarningAlertVariant,
		}

		w.WriteHeader(http.StatusMethodNotAllowed)

		if err := json.NewEncoder(w).Encode(responseBody); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	authorizationHeader := req.Header.Get("Authorization")

	if authorizationHeader == "" {
		responseBody := models.BaseResponse{
			Message:      token_constants.TokenVerifyEmptyAuthorizationHeader,
			AlertVariant: models.WarningAlertVariant,
		}

		w.WriteHeader(http.StatusUnauthorized)

		if err := json.NewEncoder(w).Encode(responseBody); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	const bearerPrefix = "Bearer "

	if !strings.HasPrefix(authorizationHeader, bearerPrefix) {
		responseBody := models.BaseResponse{
			Message:      token_constants.TokenVerifyInvalidFormat,
			AlertVariant: models.WarningAlertVariant,
		}

		w.WriteHeader(http.StatusUnauthorized)

		if err := json.NewEncoder(w).Encode(responseBody); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	token := strings.TrimPrefix(authorizationHeader, bearerPrefix)

	user, err := token_service.Verify(token)

	if err != nil {
		responseBody := models.BaseResponse{
			Message:      token_constants.TokenVerifyInvalidToken,
			AlertVariant: models.WarningAlertVariant,
		}

		w.WriteHeader(http.StatusUnauthorized)

		if err := json.NewEncoder(w).Encode(responseBody); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	responseBody := responses.TokenVerifyResponse{
		BaseResponse: models.BaseResponse{
			Message:      token_constants.TokenVerifySuccessMessage,
			AlertVariant: models.SuccessAlertVariant,
		},
		User: *user,
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
