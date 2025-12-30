package users

import "time"

type User struct {
	Id           string    `db:"id"`
	Nickname     string    `db:"nickname"`
	Email        string    `db:"email"`
	RoleId       string    `db:"role_id"`
	BannedBefore time.Time `db:"banned_before"`
	CreatedAt    time.Time `db:"created_at"`
}

type UpdateUser struct {
	Id           string    `db:"id"`
	Nickname     string    `db:"nickname"`
	Email        string    `db:"email"`
	RoleId       string    `db:"role_id"`
	BannedBefore time.Time `db:"banned_before"`
}
