package secretservice

import (
	"strings"
	"testing"

	"github.com/docker/docker-credential-helpers/credentials"
)

func TestSecretServiceHelper(t *testing.T) {
	t.Skip("test requires gnome-keyring but travis CI doesn't have it")

	creds := &credentials.Credentials{
		ServerURL: "https://foobar.docker.io:2376/v1",
		Username:  "foobar",
		Secret:    "foobarbaz",
	}

	helper := Secretservice{}
	old_auths, err := helper.List() // 1
	if err != nil {
		t.Fatal(err)
	}

	if len(old_auths) >= 1 {
		for k, v := range old_auths {
			if strings.Compare(k, creds.ServerURL) == 0 && strings.Compare(v, creds.Username) == 0 {

				if err := helper.Delete(creds.ServerURL); err != nil {
					t.Fatal(err)
				}
			}
		}

	}

	old_auths, err = helper.List()
	if err != nil {
		t.Fatal(err)
	}

	if err := helper.Add(creds); err != nil {
		t.Fatal(err)
	}

	username, secret, err := helper.Get(creds.ServerURL)
	if err != nil {
		t.Fatal(err)
	}

	if username != "foobar" {
		t.Fatalf("expected %s, got %s\n", "foobar", username)
	}

	if secret != "foobarbaz" {
		t.Fatalf("expected %s, got %s\n", "foobarbaz", secret)
	}

	new_auths, err := helper.List()
	if err != nil || (len(new_auths)-len(old_auths) != 1) {
		t.Fatal(err)
	}
	old_auths = new_auths

	if err := helper.Delete(creds.ServerURL); err != nil {
		t.Fatal(err)
	}

	new_auths, err = helper.List()
	if err != nil || (len(old_auths)-len(new_auths) != 1) {
		t.Fatal(err)
	}
	old_auths = new_auths

	helper.Add(creds)

	username, secret, err = helper.Get(creds.ServerURL)
	if err != nil {
		t.Fatal(err)
	}

	if username != "foobar" {
		t.Fatalf("expected %s, got %s\n", "foobar", username)
	}

	if secret != "foobarbaz" {
		t.Fatalf("expected %s, got %s\n", "foobarbaz", secret)
	}

	new_auths, err = helper.List()
	if (len(new_auths) - len(old_auths)) != 1 {
		t.Fatal(err)
	}
	old_auths = new_auths

	if err := helper.Delete(creds.ServerURL); err != nil {
		t.Fatal(err)
	}

	new_auths, err = helper.List()
	if (len(old_auths) - len(new_auths)) != 1 {
		t.Fatal(err)
	}
}

func TestMissingCredentials(t *testing.T) {
	t.Skip("test requires gnome-keyring but travis CI doesn't have it")

	helper := Secretservice{}
	_, _, err := helper.Get("https://adsfasdf.wrewerwer.com/asdfsdddd")
	if !credentials.IsErrCredentialsNotFound(err) {
		t.Fatalf("expected ErrCredentialsNotFound, got %v", err)
	}
}
