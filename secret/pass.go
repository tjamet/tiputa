package secret

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/afero"
)

// Pgm is the name of the pass program
const Pgm = "pass"

// SecretError is the generic error reported when failing to decode a pass secret
type SecretError struct {
	secret            []string
	passwordStorePath string
	err               error
}

func (se SecretError) Error() string {
	extra := ""
	if se.err != nil {
		extra = fmt.Sprintf(": %s", se.err.Error())
	}
	return fmt.Sprintf("failed to get secret %s in store %s%s", filepath.Join(se.secret...), se.passwordStorePath, extra)
}

// UnknownSecretError is the error returned when requesting a secret that does not exist in the pass database
type UnknownSecretError struct {
	secret            []string
	passwordStorePath string
}

func (us UnknownSecretError) Error() string {
	return fmt.Sprintf("unknown secret %s in store %s", filepath.Join(us.secret...), us.passwordStorePath)
}

// Reader defines the function an object must implement to read a secret
type Reader interface {
	Read(path ...string) (string, error)
}

// Stater allows to stat a given file
type Stater interface {
	// Stat provides an interface for os.Stat call
	Stat(path string) (os.FileInfo, error)
}

// PasswordStore olds the settings to decode secrets stored in pass
// see https://www.passwordstore.org/
type PasswordStore struct {
	// Path is the path of the password store directory
	Path string
	// Fs holds a pointer to the filesystem abstraction
	Fs Stater
	// Program is the name of the pass program
	Program string
}

// GetPath returns the path of the passwordstore directory
// it looks up at the PASSWORD_STORE_DIR environment variable and defaults
// $HOME/.password-store
func GetPath() string {
	if v, ok := os.LookupEnv("PASSWORD_STORE_DIR"); ok {
		return v
	}
	return fmt.Sprintf("%s/.password-store", homeDir())
}

// DefaultPasswordStore provide access to secrets stored in the default passwordstore
var DefaultPasswordStore = PasswordStore{
	Path:    GetPath(),
	Fs:      afero.NewOsFs(),
	Program: Pgm,
}

func (p PasswordStore) Read(path ...string) (string, error) {
	secret := filepath.Join(path...)
	secretPath := filepath.Join(p.Path, secret+".gpg")
	if _, err := p.Fs.Stat(secretPath); err != nil {
		return "", UnknownSecretError{
			secret:            path,
			passwordStorePath: p.Path,
		}
	}
	c := exec.Command(p.Program, secret)
	b := bytes.NewBuffer(nil)
	c.Stdout = b
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	err := c.Run()
	if err != nil {
		return "", SecretError{
			secret:            path,
			passwordStorePath: p.Path,
			err:               err,
		}
	}
	return b.String(), nil
}
