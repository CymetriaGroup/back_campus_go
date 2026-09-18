package mail

import (
	"context"
	"log/slog"

	"hexagonal-go-backend/internal/core/ports"
)

type LogMailer struct{ logger *slog.Logger }

func NewLogMailer(l *slog.Logger) *LogMailer { return &LogMailer{l} }
func (m *LogMailer) Send(_ context.Context, msg ports.MailMessage) error {
	m.logger.Info("email queued", "to", msg.To, "subject", msg.Subject)
	return nil
}
