package secret

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/tjamet/xgo/xtesting"
)

func TestGetPath(t *testing.T) {
	defer xtesting.InEnv("HOME", "/some/dir")()
	defer xtesting.InEnv("USERPROFILE", "/some/dir")()
	t.Run("with no environment variable, $HOME/.password-store is used", func(t *testing.T) {
		defer xtesting.NoEnv("PASSWORD_STORE_DIR")()
		assert.Equal(t, "/some/dir/.password-store", GetPath())
	})
	t.Run("with an environment variable, its value is used", func(t *testing.T) {
		defer xtesting.InEnv("PASSWORD_STORE_DIR", "./.password-store")()
		assert.Equal(t, "./.password-store", GetPath())
	})
}
func TestRead(t *testing.T) {
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/.password-store/some/folder", 0600))
	fd, err := fs.Create("/.password-store/some/folder/secret.gpg")
	assert.NoError(t, err)
	_, err = fd.WriteString("hello world")
	assert.NoError(t, err)
	store := PasswordStore{
		Fs:      fs,
		Path:    "/.password-store",
		Program: "./test",
	}
	t.Run("When secret does not exist", func(t *testing.T) {
		v, err := store.Read("some", "value")
		assert.Error(t, err)
		assert.Equal(t, "", v)
		e := err.(UnknownSecretError)
		assert.Equal(t, "/.password-store", e.passwordStorePath)
		assert.Equal(t, []string{"some", "value"}, e.secret)
		assert.Equal(t, "unknown secret some/value in store /.password-store", err.Error())
	})
	t.Run("When pass succeeds", func(t *testing.T) {
		store.Program = "./test_helpers/succeeds"
		v, err := store.Read("some", "folder", "secret")
		assert.NoError(t, err)
		assert.Equal(t, "this is the secret value", v)
	})
	t.Run("When pass fails", func(t *testing.T) {
		store.Program = "./test_helpers/fails"
		v, err := store.Read("some", "folder", "secret")
		assert.Error(t, err)
		assert.Equal(t, "", v)
		e := err.(SecretError)
		assert.Equal(t, "/.password-store", e.passwordStorePath)
		assert.Equal(t, []string{"some", "folder", "secret"}, e.secret)
		assert.Equal(t, "failed to get secret some/folder/secret in store /.password-store: exit status 1", err.Error())
	})
}
