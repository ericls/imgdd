package email

type EmailBackendType string

const (
	EmailBackendTypeSMTP EmailBackendType = "smtp"
	EmailBackendDummy    EmailBackendType = "dummy"
)

func (ebt EmailBackendType) IsValid() bool {
	switch ebt {
	case EmailBackendTypeSMTP, EmailBackendDummy:
		return true
	}
	return false
}

type EmailBackend interface {
	SendEmail(from string, to []string, subject string, htmlBody string, plainTextBody string) error
	GetFrom() string
}
