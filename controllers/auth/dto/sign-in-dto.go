package auth_dto

type SignInDTO struct {
	Email    string
	Password string
}

func (this *SignInDTO) IsInvalid() bool {
	isEmailInvalid := this.IsEmailInvalid()
	isPasswordInvalid := this.IsPasswordInvalid()

	isInvalid := isEmailInvalid || isPasswordInvalid

	return isInvalid
}

func (this *SignInDTO) IsEmailInvalid() bool {
	return len(this.Email) == 0
}

func (this *SignInDTO) IsPasswordInvalid() bool {
	return len(this.Password) == 0
}
