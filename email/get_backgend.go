package email

import "fmt"

var emailBackendCache = map[string]EmailBackend{}

func GetEmailBackendFromConfig(config *EmailConfigDef) (EmailBackend, error) {
	key := config.Hash()
	if backend, ok := emailBackendCache[key]; ok {
		return backend, nil
	}
	switch config.Type {
	case EmailBackendTypeSMTP:
		backend, err := NewSMTPBackend(config.SMTP)
		if err != nil {
			return nil, err
		}
		emailBackendCache[key] = backend
		return backend, nil
	case EmailBackendDummy:
		return NewDummyBackend(), nil
	}
	return nil, fmt.Errorf("unsupported email backend type: %s", config.Type)
}
