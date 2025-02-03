package auth_constants

const SignInMethodNotAllowedMessage string = "Método não permitido. /auth/sign-in só aceita POST."

const SignInInvalidEmailMessage string = "O email é obrigatório e tem que ser uma string."
const SignInInvalidPasswordMessage string = "A senha é obrigatória e tem que ser uma string."

const SignInDbConnectionErrorMessage string = "Erro inesperado no banco de dados ao realizar o login."

const SignInUnauthorizedMessage string = "Login ou senha inválido."
const SignInGenerateTokenErrorMessage string = "Erro inesperado no servidor ao gerar o token de autenticação."
const SignInSuccessMessage string = "Usuário autenticado com sucesso."
