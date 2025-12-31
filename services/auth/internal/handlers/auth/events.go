package auth

const (
	UserCreatedTopic = "users.registered"
)

type UserRegisteredEvent struct {
	UserId         string
	Email          string
	Code           string
	IdempotencyKey string
}
