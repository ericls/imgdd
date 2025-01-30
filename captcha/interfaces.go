package captcha

import "context"

type CaptchaClient interface {
	VerifyCaptcha(ctx context.Context, token string, action string) (bool, error)
}
