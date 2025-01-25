package email

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func buildSMTPMessage(from string, to []string, subject, htmlBody, plainTextBody string) string {
	headers := map[string]string{
		"From":    from,
		"To":      strings.Join(to, ", "),
		"Subject": subject,
	}

	var message strings.Builder

	if plainTextBody != "" && htmlBody != "" {
		boundary := uuid.New().String()
		headers["Content-Type"] = fmt.Sprintf("multipart/alternative; boundary=\"%s\"", boundary)

		for k, v := range headers {
			message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
		}
		message.WriteString("\r\n")

		message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		message.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n\r\n")
		message.WriteString(plainTextBody)
		message.WriteString("\r\n")

		message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		message.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n\r\n")
		message.WriteString(htmlBody)
		message.WriteString("\r\n")

		message.WriteString(fmt.Sprintf("--%s--", boundary))
	} else if plainTextBody != "" {
		headers["Content-Type"] = "text/plain; charset=\"utf-8\""
		for k, v := range headers {
			message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
		}
		message.WriteString("\r\n")
		message.WriteString(plainTextBody)
	} else if htmlBody != "" {
		headers["Content-Type"] = "text/html; charset=\"utf-8\""
		for k, v := range headers {
			message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
		}
		message.WriteString("\r\n")
		message.WriteString(htmlBody)
	}

	return message.String()
}
