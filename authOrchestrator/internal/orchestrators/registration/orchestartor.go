package registration

import (
	"authOrchestrator/internal/orchestrators"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/flores666/profileshare-lib/eventBus"
)

type registrationOrchestrator struct {
	logger   *slog.Logger
	producer eventBus.Producer
	consumer eventBus.Consumer
}

func NewOrchestrator(
	consumer eventBus.Consumer,
	producer eventBus.Producer,
	logger *slog.Logger,
) orchestrators.Orchestrator {
	return &registrationOrchestrator{
		logger:   logger,
		producer: producer,
		consumer: consumer,
	}
}

func (o *registrationOrchestrator) Run(ctx context.Context) error {
	//todo состояние хранить в бд (state machine)
	return o.consumer.Consume(ctx, func(data []byte) error {
		var message UserRegisteredMessage
		if err := json.Unmarshal(data, &message); err != nil {
			o.logger.Error("unmarshal error", slog.String("error", err.Error()))
			return err
		}

		return o.producer.Produce(ctx, "emails.send", getEmailMessage(message))
	})
}

func getEmailMessage(msg UserRegisteredMessage) EmailMessage {
	htmlMessage := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="ru">
<head>
  <meta charset="UTF-8">
  <title>Подтверждение email — Lumo</title>
  <style>
    body {
      font-family: Arial, sans-serif;
      background-color: #f5f6fa;
      margin: 0;
      padding: 0;
    }
    .container {
      max-width: 600px;
      margin: 40px auto;
      background-color: #ffffff;
      border-radius: 8px;
      box-shadow: 0 2px 8px rgba(0,0,0,0.05);
      padding: 20px 40px 40px 40px;
    }
    .button {
      display: inline-block;
      padding: 14px 24px;
      background-color: #18181b;
      color: #ffffff !important;
      text-decoration: none;
      border-radius: 6px;
      font-weight: bold;
    }
    .footer {
      margin-top: 30px;
      font-size: 12px;
      color: #888888;
      text-align: center;
    }
  </style>
</head>
<body>
  <div class="container">
    <h2>Подтвердите ваш email</h2>
    <p>Здравствуйте!</p>
    <p>Спасибо за регистрацию на платформе <strong>Lumo</strong>.</p>
    <p>Пожалуйста, подтвердите ваш адрес электронной почты, нажав на кнопку ниже:</p>

    <p style="text-align: center; margin: 30px 0;">
      <a href="%s" class="button">Подтвердить email</a>
    </p>

    <p>Если вы не регистрировались в Lumo, просто проигнорируйте это письмо.</p>

    <div class="footer">
      &copy; 2025 Lumo. Все права защищены.
    </div>
  </div>
</body>
</html>
`, msg.ReturnUrl)

	return EmailMessage{
		To:             msg.Email,
		Message:        htmlMessage,
		Title:          "Подтверждение регистрации",
		IdempotencyKey: msg.IdempotencyKey,
	}
}
