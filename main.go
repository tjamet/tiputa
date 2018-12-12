package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/tjamet/tiputa/secret"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientauthentication "k8s.io/client-go/pkg/apis/clientauthentication/v1beta1"
)

func ParseSecret(secret string, destination *string, reader secret.Reader) error {
	t, err := reader.Read(secret)
	if err != nil {
		return err
	}
	*destination = strings.Trim(t, "\n")
	return nil
}

type SecretParserValue struct {
	destination *string
	reader      secret.Reader
}

func (spv SecretParserValue) String() string {
	if spv.destination != nil {
		return *spv.destination
	}
	return ""
}

func (spv SecretParserValue) Set(s string) error {
	return ParseSecret(s, spv.destination, spv.reader)
}

func RegisterSecretArgs(description, prefix string, reader secret.Reader, destination *clientauthentication.ExecCredentialStatus) {
	flag.Var(SecretParserValue{&destination.Token, reader}, fmt.Sprintf("%s-token", prefix), fmt.Sprintf("the %s secret containing the kubernetes access token", description))
	flag.Var(SecretParserValue{&destination.ClientCertificateData, reader}, fmt.Sprintf("%s-client-certificate", prefix), fmt.Sprintf("the %s secret containing the kubernetes access token", description))
	flag.Var(SecretParserValue{&destination.ClientKeyData, reader}, fmt.Sprintf("%s-client-key", prefix), fmt.Sprintf("the %s secret containing the kubernetes access token", description))
}

func main() {
	status := &clientauthentication.ExecCredentialStatus{}
	RegisterSecretArgs("password store", "pass", secret.DefaultPasswordStore, status)
	flag.Parse()
	credentials := clientauthentication.ExecCredential{
		TypeMeta: metav1.TypeMeta{
			APIVersion: clientauthentication.SchemeGroupVersion.String(),
			Kind:       "ExecCredential",
		},
		Status: status,
	}
	e := json.NewEncoder(os.Stdout)
	e.Encode(credentials)
}
