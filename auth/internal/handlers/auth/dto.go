package auth

type RegisterUserRequest struct {
	Nickname  string `json:"nickname" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	ReturnUrl string `json:"returnUrl" validate:"required,url"`
}

type ConfirmUserRequest struct {
	UserId string `json:"userId" validate:"required"`
	Code   string `json:"code" validate:"required"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
