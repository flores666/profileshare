package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

type EmailsHandler interface {
	Handle([]byte) error
}

type emailsHandler struct {
	logger *slog.Logger
}

func NewEmailsHandler(logger *slog.Logger) EmailsHandler {
	return &emailsHandler{logger}
}

func (e *emailsHandler) Handle(data []byte) error {
	// todo handle

	var message EmailMessage
	if err := json.Unmarshal(data, &message); err != nil {
		e.logger.Error("unmarshal error", slog.String("error", err.Error()))
		return err
	}

	e.logger.Info(fmt.Sprintf("got message to: %s, title: %s, please handle me!", message.To, message.Title))

	return nil
}
