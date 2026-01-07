package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"mailer/internal/handlers/mailer"
	"mailer/internal/handlers/statuses"
	"mailer/internal/storage"
	"strings"
	"time"

	"github.com/flores666/profileshare-lib/utils"
)

type EmailsHandler interface {
	Handle([]byte) error
}

type emailsHandler struct {
	logger     *slog.Logger
	mailer     mailer.Service
	repository EmailsRepository
}

func NewEmailsHandler(logger *slog.Logger, mailer mailer.Service, repository EmailsRepository) EmailsHandler {
	return &emailsHandler{
		logger:     logger,
		mailer:     mailer,
		repository: repository,
	}
}

func (e *emailsHandler) Handle(data []byte) error {
	var message EmailMessage
	if err := json.Unmarshal(data, &message); err != nil {
		e.logger.Error("unmarshal error", slog.String("error", err.Error()))
		return err
	}

	idx := strings.Index(message.To, "@")
	if idx == -1 {
		e.logger.Warn("invalid email address")
		return errors.New("invalid email address")
	}

	e.logger.Info(fmt.Sprintf("got message to: %s, title: %s", message.To[:2]+"***"+message.To[idx:], message.Title))

	sentEmail, err := e.repository.GetByIdempotencyKey(context.Background(), message.IdempotencyKey)
	if err != nil {
		e.logger.Error("repository get error", slog.String("error", err.Error()))
		return err
	}

	if sentEmail != nil {
		e.logger.Warn("email already sent", slog.String("idempotencyKey", sentEmail.IdempotencyKey))
		return nil
	}

	if err = e.repository.Save(
		context.Background(),
		&storage.Email{
			Id:             utils.NewGuid(),
			Recipient:      message.To,
			Text:           message.Message,
			Subject:        message.Title,
			Status:         statuses.Sent,
			IdempotencyKey: message.IdempotencyKey,
			CreatedAt:      time.Now(),
		}); err != nil {
		e.logger.Error("repository save error", slog.String("error", err.Error()))
		return err
	}

	if err = e.mailer.Send(message.To, message.Title, message.Message); err != nil {
		e.logger.Error("send email error", slog.String("error", err.Error()))
		return err
	}

	return nil
}
