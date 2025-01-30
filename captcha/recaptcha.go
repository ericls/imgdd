package captcha

import "context"

type RecaptchaClient struct {
	serverKey string
}

func NewRecaptchaClient(secretKey string) *RecaptchaClient {
	return &RecaptchaClient{serverKey: secretKey}
}

func (r *RecaptchaClient) VerifyCaptcha(ctx context.Context, token string, action string) (bool, error) {
	return false, nil
}
