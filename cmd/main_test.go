package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tjamet/tiputa/credentials"
	clientauthentication "k8s.io/client-go/pkg/apis/clientauthentication/v1beta1"
)

func setOutput(w io.Writer) {
	output = w
}

func setOsArgs(args []string) {
	os.Args = args
}

func setGetters(g []credentials.Getter) {
	getters = g
}

type testSecretReader struct{}

func (t testSecretReader) Read(path ...string) (string, error) {
	return strings.Join(path, "/"), nil
}

func TestExecute(t *testing.T) {
	defer setOutput(output)
	defer setGetters(getters)
	defer setOsArgs(os.Args)
	setOsArgs([]string{"test", "--test-token", "hello"})
	setGetters([]credentials.Getter{})
	RegisterSecretArgs("this is a test", "test", testSecretReader{})
	b := bytes.NewBuffer(nil)
	setOutput(b)
	Execute()
	d := json.NewDecoder(b)
	c := clientauthentication.ExecCredential{}
	err := d.Decode(&c)
	assert.NoError(t, err)
	assert.Equal(t, "hello", c.Status.Token)
	assert.Equal(t, "ExecCredential", c.Kind)
	assert.Equal(t, "client.authentication.k8s.io/v1beta1", c.APIVersion)
}
