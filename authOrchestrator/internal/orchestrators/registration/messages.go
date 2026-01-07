package registration

type UserRegisteredMessage struct {
	UserId         string `json:"userId"`
	Email          string `json:"email"`
	ReturnUrl      string `json:"returnUrl"`
	IdempotencyKey string `json:"idempotencyKey"`
}

type EmailMessage struct {
	To             string `json:"to"`
	Message        string `json:"message"`
	Title          string `json:"title"`
	IdempotencyKey string `json:"idempotencyKey"`
}
