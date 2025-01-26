package identity

import (
	"net/url"
	"regexp"
	"testing"

	"github.com/ericls/imgdd/db"
	"github.com/ericls/imgdd/email"
	"github.com/ericls/imgdd/test_support"
)

func TestResetPassword(t *testing.T) {
	emailBackend := email.NewDummyBackend()
	dbConfig := TestServiceMan.GetDBConfig()
	conn := db.GetConnection(dbConfig)
	test_support.ResetDatabase(dbConfig)
	identityRepo := NewDBIdentityRepo(conn)
	orgUser, err := identityRepo.CreateUserWithOrganization("test@home.arpa", "test", "123")
	if err != nil {
		t.Fatal(err)
	}
	user := orgUser.User
	secretKey := "123"
	baseURL, err := url.Parse("https://here.home.arpa")
	if err != nil {
		t.Fatal(err)
	}
	err = SendResetPasswordEmail(emailBackend, secretKey, user, identityRepo, baseURL)
	if err != nil {
		t.Fatal(err)
	}
	recievedEmail := emailBackend.SentMessages[0].HTMLBody
	regex := regexp.MustCompile(`(?m)message=([a-zA-Z0-9._-]+)`)
	message := regex.FindStringSubmatch(recievedEmail)[1]
	t.Run("invalid message", func(t *testing.T) {
		err := ResetPassword(identityRepo, secretKey, message+"a", "newpassword")
		if err == nil {
			t.Fatal("expected error. Bad message")
		}
		if CheckPasswordByUserId(user.Id, "newpassword", identityRepo) {
			t.Fatal("password should not be updated")
		}
		if !CheckPasswordByUserId(user.Id, "123", identityRepo) {
			t.Fatal("password should not be updated")
		}
	})
	t.Run("valid message", func(t *testing.T) {
		err := ResetPassword(identityRepo, secretKey, message, "newpassword")
		if err != nil {
			t.Fatal(err)
		}
		if !CheckPasswordByUserId(user.Id, "newpassword", identityRepo) {
			t.Fatal("password not updated")
		}
	})

}
