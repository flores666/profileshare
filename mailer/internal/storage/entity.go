package storage

import "time"

type Email struct {
	Id             string    `db:"id"`
	Recipient      string    `db:"recipient"`
	Text           string    `db:"text"`
	Subject        string    `db:"subject"`
	Status         string    `db:"status"`
	IdempotencyKey string    `db:"idempotency_key"`
	CreatedAt      time.Time `db:"created_at"`
}
