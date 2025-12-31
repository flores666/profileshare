package auth

const (
	UserCreatedTopic = "users.registered"
)

type UserRegisteredMessage struct {
	UserId         string `json:"userId"`
	Email          string `json:"email"`
	Code           string `json:"code"`
	IdempotencyKey string `json:"idempotencyKey"`
}
