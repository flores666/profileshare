package registration

type UserRegisteredMessage struct {
	UserId         string `json:"userId"`
	Email          string `json:"email"`
	Code           string `json:"code"`
	IdempotencyKey string `json:"idempotencyKey"`
}

type EmailMessage struct {
	To             string `json:"to"`
	Message        string `json:"message"`
	Title          string `json:"title"`
	IdempotencyKey string `json:"idempotencyKey"`
}
