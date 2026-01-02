package handlers

type EmailMessage struct {
	To             string `json:"to"`
	Message        string `json:"message"`
	Title          string `json:"title"`
	IdempotencyKey string `json:"idempotencyKey"`
}
