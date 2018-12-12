package credentials

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testSecretReader struct{}

func (t testSecretReader) Read(path ...string) (string, error) {
	return "test---" + strings.Join(path, "--test--"), nil
}

type testSecretReaderError struct{}

func (t testSecretReaderError) Read(path ...string) (string, error) {
	return "", fmt.Errorf("test error")
}

func TestSecretGetter(t *testing.T) {
	t.Run("when only token secret is provided, token is used", func(t *testing.T) {
		status, err := Secret{
			Reader:      testSecretReader{},
			TokenSecret: "this/is/the/secret/token",
		}.Get()
		assert.NoError(t, err)
		assert.Equal(t, "test---this/is/the/secret/token", status.Token)
	})
	t.Run("when only client certificate details are provided, those details are used", func(t *testing.T) {
		status, err := Secret{
			Reader:                  testSecretReader{},
			ClientCertificateSecret: "this/is/the/secret/certificate",
			ClientKeySecret:         "this/is/the/secret/key",
		}.Get()
		assert.NoError(t, err)
		assert.Equal(t, "test---this/is/the/secret/certificate", status.ClientCertificateData)
		assert.Equal(t, "test---this/is/the/secret/key", status.ClientKeyData)
	})
	t.Run("when token secret and certificate data are provided, token is used", func(t *testing.T) {
		status, err := Secret{
			Reader:                  testSecretReader{},
			TokenSecret:             "this/is/the/secret/token",
			ClientCertificateSecret: "this/is/the/secret/certificate",
			ClientKeySecret:         "this/is/the/secret/key",
		}.Get()
		assert.NoError(t, err)
		assert.Equal(t, "test---this/is/the/secret/token", status.Token)
	})
	t.Run("when client certificate is missing but key provided, status is not found", func(t *testing.T) {
		status, err := Secret{
			Reader:                  testSecretReader{},
			ClientCertificateSecret: "this/is/the/secret/certificate",
		}.Get()
		assert.Error(t, err)
		_, ok := err.(NotProvidedError)
		assert.Truef(t, ok, "the error should be a NotProvidedError")
		assert.Nil(t, status)
	})
	t.Run("when client key is missing but certificate provided, status is not found", func(t *testing.T) {
		status, err := Secret{
			Reader:          testSecretReader{},
			ClientKeySecret: "this/is/the/secret/key",
		}.Get()
		assert.Error(t, err)
		_, ok := err.(NotProvidedError)
		assert.Truef(t, ok, "the error should be a NotProvidedError")
		assert.Nil(t, status)
	})
	t.Run("when no keys is provided, status is not found", func(t *testing.T) {
		status, err := Secret{
			Reader: testSecretReader{},
		}.Get()
		assert.Error(t, err)
		_, ok := err.(NotProvidedError)
		assert.Truef(t, ok, "the error should be a NotProvidedError")
		assert.Nil(t, status)
	})
	t.Run("when secret reading returns an error, this error is forwarded", func(t *testing.T) {
		status, err := Secret{
			Reader:      testSecretReaderError{},
			TokenSecret: "this/is/the/secret/token",
		}.Get()
		assert.Error(t, err)
		assert.Nil(t, status)
	})
	t.Run("when secret reading returns an error, this error is forwarded", func(t *testing.T) {
		status, err := Secret{
			Reader:                  testSecretReaderError{},
			ClientCertificateSecret: "this/is/the/secret/certificate",
			ClientKeySecret:         "this/is/the/secret/key",
		}.Get()
		assert.Error(t, err)
		assert.Nil(t, status)
	})
}
