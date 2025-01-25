package email

func SendEmail(backend EmailBackend, from string, to []string, subject, htmlBody, plainTextBody string) error {
	realFrom := backend.GetFrom()
	if realFrom == "" {
		realFrom = from
	}
	return backend.SendEmail(realFrom, to, subject, htmlBody, plainTextBody)
}
