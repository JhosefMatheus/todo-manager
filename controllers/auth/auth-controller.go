package authcontroller

import (
	"encoding/json"
	"io"
	"net/http"
	auth_constants "todo-manager/constants/auth"
	"todo-manager/controllers/auth/dto"
	"todo-manager/models"
	authservice "todo-manager/services/auth"
)

func SignIn(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		responseBody := models.BaseResponse{
			Message:      auth_constants.SignInMethodNotAllowedMessage,
			AlertVariant: models.WarningAlertVariant,
		}

		w.WriteHeader(http.StatusMethodNotAllowed)

		if err := json.NewEncoder(w).Encode(responseBody); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	body, err := io.ReadAll(req.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var dto dto.SignInDTO

	if err = json.Unmarshal(body, &dto); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if dto.IsInvalid() {
		message := ""

		if dto.IsEmailInvalid() {
			message += auth_constants.SignInInvalidEmailMessage
		}

		if dto.IsPasswordInvalid() {
			if len(message) > 0 {
				message += "\n"
			}

			message += auth_constants.SignInInvalidPasswordMessage
		}

		response := models.BaseResponse{
			Message:      message,
			AlertVariant: models.WarningAlertVariant,
		}

		w.WriteHeader(http.StatusForbidden)

		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	status, errResponse, signInResponse := authservice.SignIn(dto)

	w.WriteHeader(status)

	if errResponse != nil {
		if err = json.NewEncoder(w).Encode(errResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		if err = json.NewEncoder(w).Encode(signInResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
