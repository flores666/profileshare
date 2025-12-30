package users

import "time"

type CreateUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type Filter struct {
	Search string `json:"search"`
}

type UpdateUserRequest struct {
	Id           string    `db:"id"`
	Nickname     string    `db:"nickname"`
	Email        string    `db:"email"`
	RoleId       string    `db:"role_id"`
	BannedBefore time.Time `db:"banned_before"`
}
