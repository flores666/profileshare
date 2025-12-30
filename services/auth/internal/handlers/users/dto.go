package users

type CreateUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type Filter struct {
	Search string `json:"search"`
}
