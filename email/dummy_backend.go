package email

import (
	"sync"

	"github.com/ericls/imgdd/logging"
)

var dummy_backend_logger = logging.GetLogger("dummy_email_backend")

type DummyMessage struct {
	From          string
	To            []string
	Subject       string
	HTMLBody      string
	PlainTextBody string
}

type DummyBackend struct {
	SentMessages []DummyMessage
	mutex        sync.RWMutex
}

func (b *DummyBackend) SendEmail(from string, to []string, subject, htmlBody, plainTextBody string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.SentMessages = append(b.SentMessages, DummyMessage{
		From:          from,
		To:            to,
		Subject:       subject,
		HTMLBody:      htmlBody,
		PlainTextBody: plainTextBody,
	})
	dummy_backend_logger.Info().Str("from", from).Strs("to", to).Str("subject", subject).Msg("Email sent")
	return nil
}

func (b *DummyBackend) GetFrom() string {
	return "dummy@home.arpa"
}

func NewDummyBackend() *DummyBackend {
	return &DummyBackend{}
}
