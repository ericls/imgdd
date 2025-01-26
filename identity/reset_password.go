package identity

import (
	"fmt"
	"net/url"
	"time"

	dm "github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/email"
	"github.com/ericls/imgdd/signing"
)

type passwordResetMessage struct {
	UserId string `json:"userId"`
	Token  string `json:"token"`
}

func SendResetPasswordEmail(
	emailBackend email.EmailBackend,
	secretKey string,
	user *dm.User,
	identityRepo IdentityRepo,
	baseURL *url.URL,
) error {
	currentPassword := identityRepo.GetUserPassword(user.Id)
	token := makeResetPasswordToken(user.Id, currentPassword, time.Now(), secretKey)
	message, err := signing.Dumps(passwordResetMessage{
		UserId: user.Id,
		Token:  token,
	}, secretKey)
	if err != nil {
		return err
	}
	resetURL := baseURL
	resetURL.Path = "/reset_password"
	resetURL.RawQuery = url.Values{"message": {message}}.Encode()
	resetURLString := resetURL.String()
	htmlBody, err := email.RenderTemplate("reset_password.html", struct {
		ResetURL string
	}{
		ResetURL: resetURLString,
	})
	if err != nil {
		return err
	}
	return email.SendEmail(emailBackend, "", []string{user.Email}, "Password reset instruction", htmlBody, "")
}

func ResetPassword(
	identityRepo IdentityRepo,
	secretKey string,
	message string,
	newPassword string,
) error {
	var resetMessage passwordResetMessage
	err := signing.Loads(message, &resetMessage, secretKey)
	if err != nil {
		return err
	}
	if !checkResetPasswordToken(resetMessage.UserId, identityRepo.GetUserPassword(resetMessage.UserId), time.Now(), secretKey, resetMessage.Token) {
		return fmt.Errorf("invalid reset password token")
	}
	return identityRepo.UpdateUserPassword(resetMessage.UserId, newPassword)
}
