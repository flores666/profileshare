package registration

import (
	"authOrchestrator/internal/orchestrators"
	"context"
	"encoding/json"
	"eventBus"
	"fmt"
	"log/slog"
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
	htmlMessage := fmt.Sprintf(`<!DOCTYPE html>
<html lang="ru">
<head>
  <meta charset="UTF-8">
  <title>Подтвердите регистрацию на Lumo</title>
</head>
<body style="font-family: Arial, sans-serif; background-color: #f4f4f4; margin:0; padding:0;">
  <table width="100%%" cellpadding="0" cellspacing="0" style="max-width:600px; margin:40px auto; background-color:#ffffff; border-radius:8px; overflow:hidden; box-shadow:0 2px 5px rgba(0,0,0,0.1);">
    <tr>
      <td style="padding:20px; text-align:center; background-color:#4B6FFF; color:#ffffff;">
        <h1 style="margin:0; font-size:24px;">Добро пожаловать в Lumo!</h1>
      </td>
    </tr>
    <tr>
      <td style="padding:30px; color:#333333; font-size:16px; line-height:1.5;">
        <p>Здравствуйте!</p>
        <p>Спасибо за регистрацию на Lumo. Чтобы подтвердить свой аккаунт, пожалуйста, нажмите на кнопку ниже:</p>
        <p style="text-align:center; margin:30px 0;">
          <a href="%s" 
             style="display:inline-block; padding:12px 24px; background-color:#4B6FFF; color:#ffffff; text-decoration:none; border-radius:6px; font-weight:bold;">
            Подтвердить регистрацию
          </a>
        </p>
        <p>Если вы не регистрировались на Lumo, просто проигнорируйте это письмо.</p>
        <p>С уважением,<br>Команда Lumo</p>
      </td>
    </tr>
  </table>
</body>
</html>`, "https://www.lumo.com")
	//todo ссылку подставить

	return EmailMessage{
		To:             msg.Email,
		Message:        htmlMessage,
		Title:          "Подтверждение регистрации",
		IdempotencyKey: msg.IdempotencyKey,
	}
}
