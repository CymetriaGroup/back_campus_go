package mail

import (
	"context"
	"log/slog"

	"hexagonal-go-backend/internal/modules/users/application"
)

type LogMailer struct{ logger *slog.Logger }

func NewLogMailer(logger *slog.Logger) *LogMailer { return &LogMailer{logger: logger} }

func (m *LogMailer) Send(_ context.Context, message application.MailMessage) error {
	m.logger.Info("email queued", "to", message.To, "subject", message.Subject)
	return nil
}
