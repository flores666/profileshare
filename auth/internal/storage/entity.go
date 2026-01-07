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

type Token struct {
	Id              string    `db:"id"`
	UserId          string    `db:"user_id"`
	ProviderName    string    `db:"provider_name"`
	Token           string    `db:"token"`
	ExpiresAt       time.Time `db:"expires_at"`
	CreatedAt       time.Time `db:"created_at"`
	ReplacedByToken string    `db:"replaced_by_token"`
	RevokedByIp     string    `db:"revoked_by_ip"`
	RevokedAt       time.Time `db:"revoked_at"`
}
