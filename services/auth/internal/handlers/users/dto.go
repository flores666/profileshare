package users

import "time"

type CreateUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type QueryFilter struct {
	Search string `json:"search"`
}

type UpdateUserRequest struct {
	Id           string     `json:"id"`
	Nickname     *string    `json:"nickname"`
	Email        *string    `json:"email"`
	RoleId       *string    `json:"roleId"`
	BannedBefore *time.Time `json:"bannedBefore"`
}
