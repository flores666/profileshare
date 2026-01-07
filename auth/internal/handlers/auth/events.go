package auth

const (
	UserCreatedTopic = "users.registered"
)

type UserRegisteredMessage struct {
	UserId         string `json:"userId"`
	Email          string `json:"email"`
	ReturnUrl      string `json:"returnUrl"`
	IdempotencyKey string `json:"idempotencyKey"`
}
