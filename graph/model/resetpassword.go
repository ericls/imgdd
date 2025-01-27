package model

type ResetPasswordInput struct {
	Message  string `json:"message"`
	Password string `json:"password"`
}

type ResetPasswordResult struct {
	Success bool `json:"success"`
}

type SendResetPasswordEmailInput struct {
	Email string `json:"email"`
}

type SendResetPasswordEmailResult struct {
	Success bool `json:"success"`
}
