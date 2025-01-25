package email

import (
	"fmt"
	"net/smtp"

	"github.com/ericls/imgdd/logging"
)

var smtp_logger = logging.GetLogger("smtp_backend")

type SMTPBackend struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func NewSMTPBackend(config *SMTPConfigDef) (*SMTPBackend, error) {
	host := config.Host
	port := config.Port
	username := config.Username
	password := config.Password
	from := config.From

	return &SMTPBackend{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
	}, nil
}

func (s *SMTPBackend) SendEmail(from string, to []string, subject, htmlBody, plainTextBody string) error {
	if len(to) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	smtp_logger.Info().Str("from", from).Strs("to", to).Str("subject", subject).Msg("Sending email")
	message := buildSMTPMessage(from, to, subject, htmlBody, plainTextBody)

	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)
	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)
	err := smtp.SendMail(addr, auth, s.From, to, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	smtp_logger.Info().Str("from", from).Strs("to", to).Str("subject", subject).Msg("Email sent")

	return nil
}

func (s *SMTPBackend) GetFrom() string {
	return s.From
}
