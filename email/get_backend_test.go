package email_test

import (
	"testing"

	"github.com/ericls/imgdd/email"
)

func TestGetSMTPBackend(t *testing.T) {
	expected := email.SMTPBackend{
		Host:     "smtp.home.arpa",
		Port:     "587",
		Username: "user",
		Password: "pass",
		From:     "test@home.arpa",
	}
	config := email.EmailConfigDef{
		Type: email.EmailBackendTypeSMTP,
		SMTP: &email.SMTPConfigDef{
			Host:     "smtp.home.arpa",
			Port:     "587",
			Username: "user",
			Password: "pass",
			From:     "test@home.arpa",
		},
	}
	backend, err := email.GetEmailBackendFromConfig(&config)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	smtpBackend := *backend.(*email.SMTPBackend)
	if smtpBackend != expected {
		t.Errorf("Expected: %v, got: %v", expected, backend)
	}

	config = email.EmailConfigDef{
		Type: email.EmailBackendDummy,
		SMTP: &email.SMTPConfigDef{
			Host:     "smtp.home.arpa",
			Port:     "587",
			Username: "user",
			Password: "pass",
			From:     "test@home.arpa",
		},
	}
	backend, err = email.GetEmailBackendFromConfig(&config)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if _, ok := backend.(*email.DummyBackend); !ok {
		t.Errorf("Expected: %T, got: %T", email.DummyBackend{}, backend)
	}
}
