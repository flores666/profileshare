package storage

import "time"

type User struct {
	Id              string    `db:"id"`
	Nickname        string    `db:"nickname"`
	Email           string    `db:"email"`
	PasswordHash    string    `db:"password_hash"`
	Code            string    `db:"code"`
	CodeRequestedAt time.Time `db:"code_requested_at"`
	IsConfirmed     bool      `db:"is_confirmed"`
	RoleId          string    `db:"role_id"`
	BannedBefore    time.Time `db:"banned_before"`
	CreatedAt       time.Time `db:"created_at"`
}

type UpdateUser struct {
	Id       string  `db:"id"`
	Nickname *string `db:"nickname"`
	Email    *string `db:"email"`
}
