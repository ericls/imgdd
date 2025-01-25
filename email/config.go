package email

import (
	"encoding/hex"
	"hash/fnv"

	"github.com/ericls/imgdd/utils"
)

type SMTPConfigDef struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func (scd *SMTPConfigDef) Hash() string {
	h := fnv.New32()
	h.Write([]byte(scd.Host))
	h.Write([]byte(scd.Port))
	h.Write([]byte(scd.Username))
	h.Write([]byte(scd.Password))
	h.Write([]byte(scd.From))
	return hex.EncodeToString(h.Sum(nil))
}

type EmailConfigDef struct {
	Type EmailBackendType
	SMTP *SMTPConfigDef
}

func (ecd *EmailConfigDef) Hash() string {
	h := fnv.New32()
	h.Write([]byte(ecd.Type))
	if ecd.SMTP != nil {
		h.Write([]byte(ecd.SMTP.Hash()))
	}
	return hex.EncodeToString(h.Sum(nil))
}

func ReadEmailConfigFromEnv() EmailConfigDef {
	return EmailConfigDef{
		Type: EmailBackendType(utils.GetEnv("EMAIL_BACKEND_TYPE", string(EmailBackendDummy))),
		SMTP: &SMTPConfigDef{
			Host:     utils.GetEnv("EMAIL_SMTP_HOST", ""),
			Port:     utils.GetEnv("EMAIL_SMTP_PORT", ""),
			Username: utils.GetEnv("EMAIL_SMTP_USERNAME", ""),
			Password: utils.GetEnv("EMAIL_SMTP_PASSWORD", ""),
			From:     utils.GetEnv("EMAIL_SMTP_FROM", ""),
		},
	}
}
